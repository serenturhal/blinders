resource "aws_lambda_function" "dictionary" {
  function_name    = "blinders_dictionary"
  filename         = "../functions/dictionary/lambda_bundle.zip"
  handler          = "blinders.dictionary_aws_lambda_function.lambda_handler"
  source_code_hash = filebase64sha256("../functions/dictionary/lambda_bundle.zip")
  role             = aws_iam_role.lambda_role.arn
  runtime          = "python3.10"
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
}

resource "null_resource" "translate" {
  provisioner "local-exec" {
    command = "cd .. && sh ./scripts/build_golambda.sh"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}

data "archive_file" "translate" {
  depends_on = [null_resource.translate]

  type        = "zip"
  source_file = "../dist/translate"
  output_path = "../dist/translate.zip"
}

resource "aws_lambda_function" "translate" {
  function_name = "blinders_translate"
  role          = aws_iam_role.lambda_role.arn
  handler       = "translate"
  memory_size   = 128

  filename         = "../dist/translate.zip"
  source_code_hash = data.archive_file.translate.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = local.envs
  }
}
