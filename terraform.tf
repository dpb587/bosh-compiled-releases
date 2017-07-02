variable "owner" { default = "dpb587" }
variable "repository" { default = "bosh-compiled-releases" }
variable "region" { default = "us-east-1" }

provider "aws" {
  region = "${var.region}"
}

data "template_file" "envrc" {
  template = <<EOF
export AWS_DEFAULT_REGION="$${region}"
export AWS_ACCESS_KEY_ID="$${access_key_id}"
export AWS_SECRET_ACCESS_KEY="$${secret_access_key}"
EOF

  vars {
    region = "${var.region}"
    access_key_id = "${aws_iam_access_key.user.id}"
    secret_access_key = "${aws_iam_access_key.user.secret}"
  }
}

resource "null_resource" "envrc" {
  triggers {
    config_final = "${data.template_file.envrc.rendered}"
  }

  provisioner "local-exec" {
    command = "echo '${data.template_file.envrc.rendered}' > ${path.module}/.envrc.local"
  }
}

resource "aws_iam_user" "user" {
    name = "${var.repository}"
}

resource "aws_iam_access_key" "user" {
    user = "${aws_iam_user.user.name}"
}

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.owner}-${var.repository}-${var.region}"
  versioning {
    enabled = true
  }
}

data "aws_iam_policy_document" "bucket" {
  statement {
    actions = [
      "s3:GetObject",
    ]
    effect = "Allow"
    principals {
      type = "*"
      identifiers = ["*"]
    }
    resources = [
      "${aws_s3_bucket.bucket.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "bucket" {
  bucket = "${aws_s3_bucket.bucket.id}"
  policy = "${data.aws_iam_policy_document.bucket.json}"
}

data "aws_iam_policy_document" "user_s3" {
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
    ]
    effect = "Allow"
    resources = [
      "${aws_s3_bucket.bucket.arn}/*",
    ]
  }
  statement {
    actions = [
      "s3:ListBucket"
    ]
    effect = "Allow"
    resources = [
      "${aws_s3_bucket.bucket.arn}",
    ]
  }
}

resource "aws_iam_user_policy" "user_s3" {
    name = "s3"
    user = "${aws_iam_user.user.name}"
    policy = "${data.aws_iam_policy_document.user_s3.json}"
}
