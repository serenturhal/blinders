resource "aws_apigatewayv2_domain_name" "blinders" {
  domain_name = "api.peakee.co"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.blinders.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  depends_on = [aws_acm_certificate_validation.blinders]
}

resource "aws_route53_record" "blinders" {
  name    = aws_apigatewayv2_domain_name.blinders.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.blinders.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.blinders.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.blinders.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
