BIN_NAME?=app

test_cases:
	@echo "Running tests cases..."

	@echo "Good test-cases:"

	@echo "	Test-case from test task running..."
	@./$(BIN_NAME) ./cases/test_task_ex.txt
	@echo "	Test-case from test task done."

	@echo "-----"

	@echo "	Good test-case 1 is running..."
	@./$(BIN_NAME) ./cases/good1.txt
	@echo "	Good test-case 1 done."

	@echo "-----"

	@echo "	Good test-case 2 is running..."
	@./$(BIN_NAME) ./cases/good2.txt
	@echo "	Good test-case 2 done."
	@echo ""

	@echo "Arguments count validation:"

	@echo "	Test-case 1 is running..."
	@./$(BIN_NAME) ./cases/arg_count_val1.txt
	@echo "	Test-case 1 done."

	@echo "-----"

	@echo "	Test-case 2 is running..."
	@./$(BIN_NAME) ./cases/arg_count_val2.txt
	@echo "	Test-case 2 done."

	@echo "-----"

	@echo "	Test-case 3 is running..."
	@./$(BIN_NAME) ./cases/arg_count_val3.txt
	@echo "	Test-case 3 done."
	@echo ""

	@echo "Tables count validation:"

	@echo "	Test-case 1 is running..."
	@./$(BIN_NAME) ./cases/tables_count_val1.txt
	@echo "	Test-case 1 done."

	@echo "-----"

	@echo "	Test-case 2 is running..."
	@./$(BIN_NAME) ./cases/tables_count_val2.txt
	@echo "	Test-case 2 done."

	@echo "-----"

	@echo "	Test-case 3 is running..."
	@./$(BIN_NAME) ./cases/tables_count_val3.txt
	@echo "	Test-case 3 done."

	@echo "-----"

	@echo " Test-case 4 is running..."
	@./$(BIN_NAME) ./cases/tables_count_val4.txt
	@echo " Test-case 4 done."
	@echo ""

	@echo "Time validation:"

	@echo " Test-case 1 is running..."
	@./$(BIN_NAME) ./cases/time_val1.txt
	@echo " Test-case 1 done."

	@echo "-----"

	@echo " Test-case 2 is running..."
	@./$(BIN_NAME) ./cases/time_val2.txt
	@echo " Test-case 2 done."

	@echo "-----"

	@echo " Test-case 3 is running..."
	@./$(BIN_NAME) ./cases/time_val3.txt
	@echo " Test-case 3 done."

	@echo "-----"

	@echo " Test-case 4 is running..."
	@./$(BIN_NAME) ./cases/time_val4.txt
	@echo " Test-case 4 done."

	@echo "All test-cases done"