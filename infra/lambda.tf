resource "aws_lambda_function" "dictionary" {
  function_name    = "blinders-dictionary"
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
  function_name    = "blinders-translate"
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
  function_name    = "blinders-ws-connect"
  filename         = "../dist/connect.zip"
  handler          = "connect"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.connect.output_base64sha256

  environment {
    variables = local.envs
  }
}


data "archive_file" "ws_authorizer" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_dir  = "../dist/authorizer"
  output_path = "../dist/ws_authorizer.zip"
}

# unzip ws_authorizer -> handler, firebase.admin.json -> TODO: protect firebase.admin.json
resource "aws_lambda_function" "ws_authorizer" {
  function_name    = "blinders-ws-authorizer"
  filename         = "../dist/ws_authorizer.zip"
  handler          = "handler"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.ws_authorizer.output_base64sha256

  environment {
    variables = local.envs
  }
}

# use archive_file instead of pre-zip file to control source code hash (consistent with plan and apply)
data "archive_file" "disconnect" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_file = "../dist/disconnect"
  output_path = "../dist/disconnect.zip"
}

resource "aws_lambda_function" "ws_disconnect" {
  function_name    = "blinders-ws-disconnect"
  filename         = "../dist/disconnect.zip"
  handler          = "disconnect"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.disconnect.output_base64sha256


  environment {
    variables = local.envs
  }
}

# use archive_file instead of pre-zip file to control source code hash (consistent with plan and apply)
data "archive_file" "ws_chat" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_file = "../dist/wschat"
  output_path = "../dist/wschat.zip"
}

resource "aws_lambda_function" "ws_chat" {
  function_name    = "blinders-ws-chat"
  filename         = "../dist/wschat.zip"
  handler          = "wschat"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "go1.x"
  source_code_hash = data.archive_file.ws_chat.output_base64sha256


  environment {
    variables = local.envs
  }
}


data "archive_file" "rest" {
  depends_on = [null_resource.go_build]

  type        = "zip"
  source_dir  = "../dist/rest"
  output_path = "../dist/rest.zip"
}

resource "aws_lambda_function" "rest" {
  function_name    = "blinders-rest-api"
  filename         = "../dist/rest.zip"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role.arn
  depends_on       = [aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role]
  runtime          = "provided.al2" # this runtime work with our built lambda (not provided.al2023)
  architectures    = ["arm64"]
  source_code_hash = data.archive_file.rest.output_base64sha256

  environment {
    variables = local.envs
  }
}
