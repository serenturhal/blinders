resource "aws_apigatewayv2_integration" "blinders_dictionary" {
  api_id = aws_apigatewayv2_api.blinders.id

  integration_uri  = aws_lambda_function.blinders_dictionary_lambda.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_dictionary" {
  api_id = aws_apigatewayv2_api.blinders.id

  route_key = "GET /dictionary"
  target    = "integrations/${aws_apigatewayv2_integration.blinders_dictionary.id}"
}

resource "aws_lambda_permission" "blinders" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.blinders_dictionary_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.blinders.execution_arn}/*/*"
}

output "get_dictionary_api" {
  value = "https://${aws_apigatewayv2_api_mapping.blinders_v1.domain_name}/${aws_apigatewayv2_api_mapping.blinders_v1.api_mapping_key}/dictionary"
}
