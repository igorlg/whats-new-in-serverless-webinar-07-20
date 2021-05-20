.PHONY: build build-template deploy
.PHONY: deploy-autoscaling deploy-sam
.PHONY: test clean purge status

AWS_REGION_NAME			= ap-southeast-2
FUNCTION_STACK_NAME		= whatsnew-lambda-pc-main
AUTOSCALING_STACK_NAME		= whatsnew-lambda-pc-autoscaling

FUNC_TEMPLATE_FILE		= template.yaml
FUNC_TEMPLATE_FILE_TMP		= template.tmp.yaml
AUTOSC_TEMPLATE_FILE		= autoscaling.yaml

TEST_CONFIG_FILE		= test/config.yaml
TEST_TEMPLATE_FILE		= test/config_template.yaml

FUNCTION_NAME			= "$(shell aws ssm get-parameter --name '/$(FUNCTION_STACK_NAME)/function/name'  | jq -r '.Parameter.Value' 2>/dev/null)"
FUNCTION_ALIAS			= "$(shell aws ssm get-parameter --name '/$(FUNCTION_STACK_NAME)/function/alias' | jq -r '.Parameter.Value' 2>/dev/null)"

build: clean
	sed -e "s/DASHBOARD_BODY/$$(jq -c -r -M '.' dashboard.json | sed 's/\//\\\//g')/" \
		$(FUNC_TEMPLATE_FILE) > $(FUNC_TEMPLATE_FILE_TMP)
	sam build \
		-t $(FUNC_TEMPLATE_FILE_TMP) \
		$(SAM_BUILD_OPTS)


deploy: deploy-sam deploy-autoscaling

deploy-sam: build
	sam deploy \
		--stack-name $(FUNCTION_STACK_NAME) \
		$(SAM_DEPLOY_OPTS)

deploy-autoscaling:
	aws cloudformation deploy $(DEPLOY_AUTOSCALING_OPTS) \
		--stack-name $(AUTOSCALING_STACK_NAME) \
		--parameter-overrides FunctionStackName=$(FUNCTION_STACK_NAME) \
		--template-file $(AUTOSC_TEMPLATE_FILE)


test:
	sed	-e "s/<FUNCTION_ALIAS>/$(FUNCTION_ALIAS)/g" \
		-e "s/<FUNCTION_NAME>/$(FUNCTION_NAME)/g" \
		-e "s/<AWS_REGION>/$(AWS_REGION_NAME)/g" \
		$(TEST_TEMPLATE_FILE) > $(TEST_CONFIG_FILE)

	AWS_REGION=$(AWS_REGION_NAME) \
		artillery run $(TEST_CONFIG_FILE)

	# @rm -f $(TEST_CONFIG_FILE)


clean:
	rm -rf .aws-sam/ $(FUNC_TEMPLATE_FILE_TMP) $(TEST_CONFIG_FILE)

purge: clean
	aws cloudformation delete-stack --stack-name $(FUNCTION_STACK_NAME)
	aws cloudformation delete-stack --stack-name $(AUTOSCALING_STACK_NAME)

status:
	@echo "Status of stacks:"
	@echo -n "$(FUNCTION_STACK_NAME):	     "
	@aws cloudformation describe-stacks --stack-name $(FUNCTION_STACK_NAME) | jq -r '.Stacks[0].StackStatus'
	@echo -n "$(AUTOSCALING_STACK_NAME):	     "
	@aws cloudformation describe-stacks --stack-name $(AUTOSCALING_STACK_NAME) | jq -r '.Stacks[0].StackStatus'

