.PHONY: build build-template deploy deploy-autoscaling deploy-sam test clean purge

AWS_REGION_NAME			= ap-southeast-2
FUNCTION_STACK_NAME    	= whatsnew-lambda-pc-main
AUTOSCALING_STACK_NAME 	= whatsnew-lambda-pc-autoscaling

FUNCTION_NAME 			= "$(shell aws ssm get-parameter --name '/$(FUNCTION_STACK_NAME)/function/name'  | jq -r '.Parameter.Value' 2>/dev/null)"
FUNCTION_ALIAS 			= "$(shell aws ssm get-parameter --name '/$(FUNCTION_STACK_NAME)/function/alias' | jq -r '.Parameter.Value' 2>/dev/null)"

TEST_CONFIG_FILE		= test/config.yaml
TEST_TEMPLATE_FILE  	= test/config_template.yaml

build: clean
	sam build -u

deploy: deploy-sam deploy-autoscaling

deploy-autoscaling:
	aws cloudformation deploy \
		--stack-name $(AUTOSCALING_STACK_NAME) \
		--parameter-overrides FunctionStackName=$(FUNCTION_STACK_NAME) \
		--template-file autoscaling.yaml \
	|| true

deploy-sam: build
	sam deploy \
		--stack-name $(FUNCTION_STACK_NAME) \
	|| true

test:
	@sed    -e "s/<FUNCTION_ALIAS>/$(FUNCTION_ALIAS)/g" \
			-e "s/<FUNCTION_NAME>/$(FUNCTION_NAME)/g" \
			-e "s/<AWS_REGION>/$(AWS_REGION_NAME)/g" \
			$(TEST_TEMPLATE_FILE) > $(TEST_CONFIG_FILE)

	AWS_REGION=$(AWS_REGION_NAME) \
	artillery run $(TEST_CONFIG_FILE)

	@rm -f $(TEST_CONFIG_FILE)

clean:
	rm -rf .aws-sam/
	rm -rf test/config.yaml

purge: clean
	aws cloudformation delete-stack --stack-name $(FUNCTION_STACK_NAME)
	aws cloudformation delete-stack --stack-name $(AUTOSCALING_STACK_NAME)
