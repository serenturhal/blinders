resource "aws_lambda_function" "dictionary" {
  function_name    = "blinders_dictionary"
  filename         = "../functions/dictionary/lambda_bundle.zip"
  handler          = "blinders.dictionary_aws_lambda_function.lambda_handler"
  source_code_hash = filebase64sha256("../functions/dictionary/lambda_bundle.zip")
  role             = aws_iam_role.lambda_role.arn
  runtime          = "python3.10"
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
}

resource "null_resource" "go_build" {
  provisioner "local-exec" {
    command = "cd .. && sh ./scripts/build_golambda.sh"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}

# use archive_file instead of pre-zip file to control source code hash (consistent with plan and apply)
data "archive_file" "translate" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_file = "../dist/translate"
  output_path = "../dist/translate.zip"
}

resource "aws_lambda_function" "translate" {
  function_name    = "blinders_translate"
  filename         = "../dist/translate.zip"
  handler          = "translate"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  memory_size      = 128
  runtime          = "go1.x"
  source_code_hash = data.archive_file.translate.output_base64sha256

  environment {
    variables = local.envs
  }
}


# use archive_file instead of pre-zip file to control source code hash (consistent with plan and apply)
data "archive_file" "connect" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_file = "../dist/connect"
  output_path = "../dist/connect.zip"
}

resource "aws_lambda_function" "ws_connect" {
  function_name    = "blinders_ws_connect"
  filename         = "../dist/connect.zip"
  handler          = "connect"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.connect.output_base64sha256
}

# use archive_file instead of pre-zip file to control source code hash (consistent with plan and apply)
data "archive_file" "disconnect" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_file = "../dist/disconnect"
  output_path = "../dist/disconnect.zip"
}

resource "aws_lambda_function" "ws_disconnect" {
  function_name    = "blinders_ws_disconnect"
  filename         = "../dist/disconnect.zip"
  handler          = "disconnect"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.disconnect.output_base64sha256
}
