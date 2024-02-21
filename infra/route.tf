# dictionary route
resource "aws_apigatewayv2_integration" "dictionary" {
  api_id = aws_apigatewayv2_api.http_api.id

  integration_uri  = aws_lambda_function.dictionary.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_dictionary" {
  api_id = aws_apigatewayv2_api.http_api.id

  route_key = "GET /dictionary"
  target    = "integrations/${aws_apigatewayv2_integration.dictionary.id}"
}

resource "aws_lambda_permission" "dictionary" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dictionary.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_dictionary_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/dictionary"
}

# translate route
resource "aws_apigatewayv2_integration" "translate" {
  api_id = aws_apigatewayv2_api.http_api.id

  integration_uri  = aws_lambda_function.translate.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "get_translate" {
  api_id = aws_apigatewayv2_api.http_api.id

  route_key = "GET /translate"
  target    = "integrations/${aws_apigatewayv2_integration.translate.id}"
}

resource "aws_lambda_permission" "translate" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.translate.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "get_translate_api" {
  value = "https://${aws_apigatewayv2_api_mapping.http_api_v1.domain_name}/${aws_apigatewayv2_api_mapping.http_api_v1.api_mapping_key}/translate"
}
