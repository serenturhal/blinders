resource "aws_lambda_function" "blinders_dictionary_lambda" {
  runtime          = "python3.10"
  filename         = "../functions/dictionary/lambda_bundle.zip"
  function_name    = "Blinders_Dictionary_Lambda_Function"
  handler          = "blinders.dictionary_aws_lambda_function.lambda_handler"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  source_code_hash = filebase64sha256("../functions/dictionary/lambda_bundle.zip")
}
