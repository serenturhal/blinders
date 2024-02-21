resource "aws_lambda_function" "dictionary" {
  runtime          = "python3.10"
  filename         = "../functions/dictionary/lambda_bundle.zip"
  function_name    = "Blinders_Dictionary_Lambda_Function"
  handler          = "blinders.dictionary_aws_lambda_function.lambda_handler"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../functions/dictionary/lambda_bundle.zip")
}

resource "null_resource" "translate" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ../dist/ ../functions/translate/"
  }

  # triggers = {
  #   always_run = "${timestamp()}"
  # }
}

data "archive_file" "translate" {
  depends_on = [null_resource.translate]

  type        = "zip"
  source_file = "../dist/translate"
  output_path = "../dist/translate.zip"
}

resource "aws_lambda_function" "translate" {
  function_name = "Blinders_translate_Lambda_Function"
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
