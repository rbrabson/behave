=== RUN   TestStatus_String
--- PASS: TestStatus_String (0.00s)
=== RUN   TestNew
--- PASS: TestNew (0.00s)
=== RUN   TestBehaviorTree_Tick
--- PASS: TestBehaviorTree_Tick (0.00s)
=== RUN   TestBehaviorTree_String
--- PASS: TestBehaviorTree_String (0.00s)
=== RUN   TestAction_Tick
=== RUN   TestAction_Tick/no_run_func
=== RUN   TestAction_Tick/success
=== RUN   TestAction_Tick/failure
=== RUN   TestAction_Tick/running
--- PASS: TestAction_Tick (0.00s)
    --- PASS: TestAction_Tick/no_run_func (0.00s)
    --- PASS: TestAction_Tick/success (0.00s)
    --- PASS: TestAction_Tick/failure (0.00s)
    --- PASS: TestAction_Tick/running (0.00s)
=== RUN   TestCondition_Tick
=== RUN   TestCondition_Tick/no_check_func
=== RUN   TestCondition_Tick/success
=== RUN   TestCondition_Tick/failure
--- PASS: TestCondition_Tick (0.00s)
    --- PASS: TestCondition_Tick/no_check_func (0.00s)
    --- PASS: TestCondition_Tick/success (0.00s)
    --- PASS: TestCondition_Tick/failure (0.00s)
=== RUN   TestSequence_Tick
=== RUN   TestSequence_Tick/empty_sequence
=== RUN   TestSequence_Tick/all_success
=== RUN   TestSequence_Tick/first_fails
=== RUN   TestSequence_Tick/second_fails
=== RUN   TestSequence_Tick/first_running
--- PASS: TestSequence_Tick (0.00s)
    --- PASS: TestSequence_Tick/empty_sequence (0.00s)
    --- PASS: TestSequence_Tick/all_success (0.00s)
    --- PASS: TestSequence_Tick/first_fails (0.00s)
    --- PASS: TestSequence_Tick/second_fails (0.00s)
    --- PASS: TestSequence_Tick/first_running (0.00s)
=== RUN   TestSelector_Tick
=== RUN   TestSelector_Tick/empty_selector
=== RUN   TestSelector_Tick/first_succeeds
=== RUN   TestSelector_Tick/second_succeeds
=== RUN   TestSelector_Tick/all_fail
=== RUN   TestSelector_Tick/first_running
--- PASS: TestSelector_Tick (0.00s)
    --- PASS: TestSelector_Tick/empty_selector (0.00s)
    --- PASS: TestSelector_Tick/first_succeeds (0.00s)
    --- PASS: TestSelector_Tick/second_succeeds (0.00s)
    --- PASS: TestSelector_Tick/all_fail (0.00s)
    --- PASS: TestSelector_Tick/first_running (0.00s)
=== RUN   TestComplexBehaviorTree
--- PASS: TestComplexBehaviorTree (0.00s)
=== RUN   TestComposite_Tick
=== RUN   TestComposite_Tick/no_conditions_or_child
=== RUN   TestComposite_Tick/no_conditions,_child_succeeds
=== RUN   TestComposite_Tick/no_conditions,_child_fails
=== RUN   TestComposite_Tick/single_condition_succeeds,_no_child
=== RUN   TestComposite_Tick/single_condition_succeeds,_child_succeeds
=== RUN   TestComposite_Tick/single_condition_succeeds,_child_fails
=== RUN   TestComposite_Tick/single_condition_fails,_child_not_executed
=== RUN   TestComposite_Tick/single_condition_running
=== RUN   TestComposite_Tick/multiple_conditions_all_succeed,_child_succeeds
=== RUN   TestComposite_Tick/multiple_conditions,_first_fails
=== RUN   TestComposite_Tick/multiple_conditions,_second_fails
=== RUN   TestComposite_Tick/multiple_conditions,_first_running
--- PASS: TestComposite_Tick (0.00s)
    --- PASS: TestComposite_Tick/no_conditions_or_child (0.00s)
    --- PASS: TestComposite_Tick/no_conditions,_child_succeeds (0.00s)
    --- PASS: TestComposite_Tick/no_conditions,_child_fails (0.00s)
    --- PASS: TestComposite_Tick/single_condition_succeeds,_no_child (0.00s)
    --- PASS: TestComposite_Tick/single_condition_succeeds,_child_succeeds (0.00s)
    --- PASS: TestComposite_Tick/single_condition_succeeds,_child_fails (0.00s)
    --- PASS: TestComposite_Tick/single_condition_fails,_child_not_executed (0.00s)
    --- PASS: TestComposite_Tick/single_condition_running (0.00s)
    --- PASS: TestComposite_Tick/multiple_conditions_all_succeed,_child_succeeds (0.00s)
    --- PASS: TestComposite_Tick/multiple_conditions,_first_fails (0.00s)
    --- PASS: TestComposite_Tick/multiple_conditions,_second_fails (0.00s)
    --- PASS: TestComposite_Tick/multiple_conditions,_first_running (0.00s)
=== RUN   TestComposite_Reset
--- PASS: TestComposite_Reset (0.00s)
=== RUN   TestParallel_Tick
=== RUN   TestParallel_Tick/empty_parallel
=== RUN   TestParallel_Tick/all_succeed,_need_all
=== RUN   TestParallel_Tick/all_succeed,_need_one
=== RUN   TestParallel_Tick/one_succeeds,_need_one
=== RUN   TestParallel_Tick/one_succeeds,_need_two
=== RUN   TestParallel_Tick/all_fail
=== RUN   TestParallel_Tick/one_running,_one_success,_need_two
=== RUN   TestParallel_Tick/one_running,_one_failure,_need_two
--- PASS: TestParallel_Tick (0.00s)
    --- PASS: TestParallel_Tick/empty_parallel (0.00s)
    --- PASS: TestParallel_Tick/all_succeed,_need_all (0.00s)
    --- PASS: TestParallel_Tick/all_succeed,_need_one (0.00s)
    --- PASS: TestParallel_Tick/one_succeeds,_need_one (0.00s)
    --- PASS: TestParallel_Tick/one_succeeds,_need_two (0.00s)
    --- PASS: TestParallel_Tick/all_fail (0.00s)
    --- PASS: TestParallel_Tick/one_running,_one_success,_need_two (0.00s)
    --- PASS: TestParallel_Tick/one_running,_one_failure,_need_two (0.00s)
=== RUN   TestParallel_MinSuccessCount_Validation
=== RUN   TestParallel_MinSuccessCount_Validation/zero_min_success_count,_should_be_1
=== RUN   TestParallel_MinSuccessCount_Validation/min_success_count_greater_than_children,_should_be_clamped
--- PASS: TestParallel_MinSuccessCount_Validation (0.00s)
    --- PASS: TestParallel_MinSuccessCount_Validation/zero_min_success_count,_should_be_1 (0.00s)
    --- PASS: TestParallel_MinSuccessCount_Validation/min_success_count_greater_than_children,_should_be_clamped (0.00s)
=== RUN   TestParallel_Reset
--- PASS: TestParallel_Reset (0.00s)
=== RUN   TestAction_Reset
--- PASS: TestAction_Reset (0.00s)
=== RUN   TestCondition_Reset
--- PASS: TestCondition_Reset (0.00s)
=== RUN   TestSequence_Reset
--- PASS: TestSequence_Reset (0.00s)
=== RUN   TestSelector_Reset
--- PASS: TestSelector_Reset (0.00s)
=== RUN   TestBehaviorTree_Reset
--- PASS: TestBehaviorTree_Reset (0.00s)
=== RUN   TestComplexBehaviorTreeWithNewNodes
--- PASS: TestComplexBehaviorTreeWithNewNodes (0.00s)
=== RUN   TestRetry_Tick
=== RUN   TestRetry_Tick/no_child
=== RUN   TestRetry_Tick/child_succeeds_immediately
=== RUN   TestRetry_Tick/child_running
--- PASS: TestRetry_Tick (0.00s)
    --- PASS: TestRetry_Tick/no_child (0.00s)
    --- PASS: TestRetry_Tick/child_succeeds_immediately (0.00s)
    --- PASS: TestRetry_Tick/child_running (0.00s)
=== RUN   TestRetry_FailureRetry
--- PASS: TestRetry_FailureRetry (0.00s)
=== RUN   TestRetry_ChildReset
--- PASS: TestRetry_ChildReset (0.00s)
=== RUN   TestRetry_Reset
--- PASS: TestRetry_Reset (0.00s)
=== RUN   TestRetry_String
=== RUN   TestRetry_String/no_child
=== RUN   TestRetry_String/with_child
--- PASS: TestRetry_String (0.00s)
    --- PASS: TestRetry_String/no_child (0.00s)
    --- PASS: TestRetry_String/with_child (0.00s)
=== RUN   TestRepeat_Tick
=== RUN   TestRepeat_Tick/no_child
=== RUN   TestRepeat_Tick/child_succeeds_first_time
=== RUN   TestRepeat_Tick/child_fails
=== RUN   TestRepeat_Tick/child_running
--- PASS: TestRepeat_Tick (0.00s)
    --- PASS: TestRepeat_Tick/no_child (0.00s)
    --- PASS: TestRepeat_Tick/child_succeeds_first_time (0.00s)
    --- PASS: TestRepeat_Tick/child_fails (0.00s)
    --- PASS: TestRepeat_Tick/child_running (0.00s)
=== RUN   TestRepeat_RepeatsUntilFailure
--- PASS: TestRepeat_RepeatsUntilFailure (0.00s)
=== RUN   TestRepeat_ChildReset
--- PASS: TestRepeat_ChildReset (0.00s)
=== RUN   TestRepeat_Reset
--- PASS: TestRepeat_Reset (0.00s)
=== RUN   TestRepeat_String
=== RUN   TestRepeat_String/no_child
=== RUN   TestRepeat_String/with_child
--- PASS: TestRepeat_String (0.00s)
    --- PASS: TestRepeat_String/no_child (0.00s)
    --- PASS: TestRepeat_String/with_child (0.00s)
=== RUN   TestComplexBehaviorTreeWithRepeat
--- PASS: TestComplexBehaviorTreeWithRepeat (0.00s)
=== RUN   TestInvert_Tick
=== RUN   TestInvert_Tick/no_child
=== RUN   TestInvert_Tick/child_succeeds,_invert_to_failure
=== RUN   TestInvert_Tick/child_fails,_invert_to_success
=== RUN   TestInvert_Tick/child_running,_pass_through
=== RUN   TestInvert_Tick/child_ready,_pass_through
--- PASS: TestInvert_Tick (0.00s)
    --- PASS: TestInvert_Tick/no_child (0.00s)
    --- PASS: TestInvert_Tick/child_succeeds,_invert_to_failure (0.00s)
    --- PASS: TestInvert_Tick/child_fails,_invert_to_success (0.00s)
    --- PASS: TestInvert_Tick/child_running,_pass_through (0.00s)
    --- PASS: TestInvert_Tick/child_ready,_pass_through (0.00s)
=== RUN   TestInvert_InversionBehavior
--- PASS: TestInvert_InversionBehavior (0.00s)
=== RUN   TestInvert_Reset
--- PASS: TestInvert_Reset (0.00s)
=== RUN   TestInvert_String
=== RUN   TestInvert_String/no_child
=== RUN   TestInvert_String/with_child
--- PASS: TestInvert_String (0.00s)
    --- PASS: TestInvert_String/no_child (0.00s)
    --- PASS: TestInvert_String/with_child (0.00s)
=== RUN   TestComplexBehaviorTreeWithInvert
--- PASS: TestComplexBehaviorTreeWithInvert (0.00s)
=== RUN   TestAlwaysSuccess_Tick
=== RUN   TestAlwaysSuccess_Tick/no_child
=== RUN   TestAlwaysSuccess_Tick/child_succeeds,_return_success
=== RUN   TestAlwaysSuccess_Tick/child_fails,_still_return_success
=== RUN   TestAlwaysSuccess_Tick/child_running,_still_return_success
=== RUN   TestAlwaysSuccess_Tick/child_ready,_still_return_success
--- PASS: TestAlwaysSuccess_Tick (0.00s)
    --- PASS: TestAlwaysSuccess_Tick/no_child (0.00s)
    --- PASS: TestAlwaysSuccess_Tick/child_succeeds,_return_success (0.00s)
    --- PASS: TestAlwaysSuccess_Tick/child_fails,_still_return_success (0.00s)
    --- PASS: TestAlwaysSuccess_Tick/child_running,_still_return_success (0.00s)
    --- PASS: TestAlwaysSuccess_Tick/child_ready,_still_return_success (0.00s)
=== RUN   TestAlwaysSuccess_ChildExecution
--- PASS: TestAlwaysSuccess_ChildExecution (0.00s)
=== RUN   TestAlwaysSuccess_Reset
--- PASS: TestAlwaysSuccess_Reset (0.00s)
=== RUN   TestAlwaysSuccess_String
=== RUN   TestAlwaysSuccess_String/no_child
=== RUN   TestAlwaysSuccess_String/with_child
--- PASS: TestAlwaysSuccess_String (0.00s)
    --- PASS: TestAlwaysSuccess_String/no_child (0.00s)
    --- PASS: TestAlwaysSuccess_String/with_child (0.00s)
=== RUN   TestComplexBehaviorTreeWithAlwaysSuccess
--- PASS: TestComplexBehaviorTreeWithAlwaysSuccess (0.00s)
=== RUN   TestAlwaysFailure_Tick
=== RUN   TestAlwaysFailure_Tick/no_child
=== RUN   TestAlwaysFailure_Tick/child_succeeds,_still_return_failure
=== RUN   TestAlwaysFailure_Tick/child_fails,_return_failure
=== RUN   TestAlwaysFailure_Tick/child_running,_still_return_failure
=== RUN   TestAlwaysFailure_Tick/child_ready,_still_return_failure
--- PASS: TestAlwaysFailure_Tick (0.00s)
    --- PASS: TestAlwaysFailure_Tick/no_child (0.00s)
    --- PASS: TestAlwaysFailure_Tick/child_succeeds,_still_return_failure (0.00s)
    --- PASS: TestAlwaysFailure_Tick/child_fails,_return_failure (0.00s)
    --- PASS: TestAlwaysFailure_Tick/child_running,_still_return_failure (0.00s)
    --- PASS: TestAlwaysFailure_Tick/child_ready,_still_return_failure (0.00s)
=== RUN   TestAlwaysFailure_ChildExecution
--- PASS: TestAlwaysFailure_ChildExecution (0.00s)
=== RUN   TestAlwaysFailure_Reset
--- PASS: TestAlwaysFailure_Reset (0.00s)
=== RUN   TestAlwaysFailure_String
=== RUN   TestAlwaysFailure_String/no_child
=== RUN   TestAlwaysFailure_String/with_child
--- PASS: TestAlwaysFailure_String (0.00s)
    --- PASS: TestAlwaysFailure_String/no_child (0.00s)
    --- PASS: TestAlwaysFailure_String/with_child (0.00s)
=== RUN   TestComplexBehaviorTreeWithAlwaysFailure
--- PASS: TestComplexBehaviorTreeWithAlwaysFailure (0.00s)
=== RUN   TestAlwaysSuccessAndAlwaysFailureTogether
--- PASS: TestAlwaysSuccessAndAlwaysFailureTogether (0.00s)
=== RUN   TestRepeatN_Tick
=== RUN   TestRepeatN_Tick/no_child
=== RUN   TestRepeatN_Tick/max_count_zero
=== RUN   TestRepeatN_Tick/single_execution
=== RUN   TestRepeatN_Tick/multiple_executions_with_success
=== RUN   TestRepeatN_Tick/multiple_executions_with_final_failure
--- PASS: TestRepeatN_Tick (0.00s)
    --- PASS: TestRepeatN_Tick/no_child (0.00s)
    --- PASS: TestRepeatN_Tick/max_count_zero (0.00s)
    --- PASS: TestRepeatN_Tick/single_execution (0.00s)
    --- PASS: TestRepeatN_Tick/multiple_executions_with_success (0.00s)
    --- PASS: TestRepeatN_Tick/multiple_executions_with_final_failure (0.00s)
=== RUN   TestRepeatN_ChildReset
--- PASS: TestRepeatN_ChildReset (0.00s)
=== RUN   TestRepeatN_Reset
--- PASS: TestRepeatN_Reset (0.00s)
=== RUN   TestRepeatN_String
=== RUN   TestRepeatN_String/no_child
=== RUN   TestRepeatN_String/with_child,_partial_execution
--- PASS: TestRepeatN_String (0.00s)
    --- PASS: TestRepeatN_String/no_child (0.00s)
    --- PASS: TestRepeatN_String/with_child,_partial_execution (0.00s)
=== RUN   TestRepeatN_RunningChild
--- PASS: TestRepeatN_RunningChild (0.00s)
=== RUN   TestComplexBehaviorTreeWithRepeatN
--- PASS: TestComplexBehaviorTreeWithRepeatN (0.00s)
=== RUN   TestWhileSuccess_Tick
=== RUN   TestWhileSuccess_Tick/no_child
=== RUN   TestWhileSuccess_Tick/child_succeeds_continuously
=== RUN   TestWhileSuccess_Tick/child_fails_after_success
=== RUN   TestWhileSuccess_Tick/child_running
--- PASS: TestWhileSuccess_Tick (0.00s)
    --- PASS: TestWhileSuccess_Tick/no_child (0.00s)
    --- PASS: TestWhileSuccess_Tick/child_succeeds_continuously (0.00s)
    --- PASS: TestWhileSuccess_Tick/child_fails_after_success (0.00s)
    --- PASS: TestWhileSuccess_Tick/child_running (0.00s)
=== RUN   TestWhileSuccess_ChildReset
--- PASS: TestWhileSuccess_ChildReset (0.00s)
=== RUN   TestWhileSuccess_Reset
--- PASS: TestWhileSuccess_Reset (0.00s)
=== RUN   TestWhileSuccess_String
=== RUN   TestWhileSuccess_String/no_child
=== RUN   TestWhileSuccess_String/with_child
--- PASS: TestWhileSuccess_String (0.00s)
    --- PASS: TestWhileSuccess_String/no_child (0.00s)
    --- PASS: TestWhileSuccess_String/with_child (0.00s)
=== RUN   TestWhileSuccess_MixedResults
--- PASS: TestWhileSuccess_MixedResults (0.00s)
=== RUN   TestComplexBehaviorTreeWithWhileSuccess
--- PASS: TestComplexBehaviorTreeWithWhileSuccess (0.00s)
=== RUN   TestWhileFailure_Tick
=== RUN   TestWhileFailure_Tick/no_child
=== RUN   TestWhileFailure_Tick/child_succeeds_immediately
=== RUN   TestWhileFailure_Tick/child_fails
=== RUN   TestWhileFailure_Tick/child_running
--- PASS: TestWhileFailure_Tick (0.00s)
    --- PASS: TestWhileFailure_Tick/no_child (0.00s)
    --- PASS: TestWhileFailure_Tick/child_succeeds_immediately (0.00s)
    --- PASS: TestWhileFailure_Tick/child_fails (0.00s)
    --- PASS: TestWhileFailure_Tick/child_running (0.00s)
=== RUN   TestWhileFailure_ChildReset
--- PASS: TestWhileFailure_ChildReset (0.00s)
=== RUN   TestWhileFailure_Reset
--- PASS: TestWhileFailure_Reset (0.00s)
=== RUN   TestWhileFailure_String
=== RUN   TestWhileFailure_String/no_child
=== RUN   TestWhileFailure_String/with_child,_running_state
--- PASS: TestWhileFailure_String (0.00s)
    --- PASS: TestWhileFailure_String/no_child (0.00s)
    --- PASS: TestWhileFailure_String/with_child,_running_state (0.00s)
=== RUN   TestWhileFailure_MixedResults
--- PASS: TestWhileFailure_MixedResults (0.00s)
=== RUN   TestComplexBehaviorTreeWithWhileFailure
--- PASS: TestComplexBehaviorTreeWithWhileFailure (0.00s)
=== RUN   TestLog_Tick
=== RUN   TestLog_Tick/no_child
2025/09/04 16:38:50 WARN Log node has no child status=Failure
=== RUN   TestLog_Tick/child_succeeds
2025/09/04 16:38:50 INFO Success test child_status=Success child_type=Action
=== RUN   TestLog_Tick/child_fails
2025/09/04 16:38:50 WARN Failure test child_status=Failure child_type=Action
=== RUN   TestLog_Tick/child_running
=== RUN   TestLog_Tick/empty_message
2025/09/04 16:38:50 INFO Log node executed child_status=Success child_type=Action
--- PASS: TestLog_Tick (0.00s)
    --- PASS: TestLog_Tick/no_child (0.00s)
    --- PASS: TestLog_Tick/child_succeeds (0.00s)
    --- PASS: TestLog_Tick/child_fails (0.00s)
    --- PASS: TestLog_Tick/child_running (0.00s)
    --- PASS: TestLog_Tick/empty_message (0.00s)
=== RUN   TestLog_ChildExecution
2025/09/04 16:38:50 INFO Execution test child_status=Success child_type=Action
2025/09/04 16:38:50 INFO Execution test child_status=Success child_type=Action
--- PASS: TestLog_ChildExecution (0.00s)
=== RUN   TestLog_Reset
2025/09/04 16:38:50 INFO Reset test child_status=Success child_type=Action
--- PASS: TestLog_Reset (0.00s)
=== RUN   TestLog_String
=== RUN   TestLog_String/no_child
=== RUN   TestLog_String/with_child_and_message
2025/09/04 16:38:50 INFO Action log child_status=Success child_type=Action
=== RUN   TestLog_String/empty_message
2025/09/04 16:38:50 INFO Log node executed child_status=Success child_type=Action
--- PASS: TestLog_String (0.00s)
    --- PASS: TestLog_String/no_child (0.00s)
    --- PASS: TestLog_String/with_child_and_message (0.00s)
    --- PASS: TestLog_String/empty_message (0.00s)
=== RUN   TestLog_GetChildType
=== RUN   TestLog_GetChildType/nil_child
=== RUN   TestLog_GetChildType/Action
=== RUN   TestLog_GetChildType/Condition
=== RUN   TestLog_GetChildType/Sequence
=== RUN   TestLog_GetChildType/Selector
=== RUN   TestLog_GetChildType/Parallel
=== RUN   TestLog_GetChildType/Composite
=== RUN   TestLog_GetChildType/Retry
=== RUN   TestLog_GetChildType/Repeat
=== RUN   TestLog_GetChildType/RepeatN
=== RUN   TestLog_GetChildType/Invert
=== RUN   TestLog_GetChildType/AlwaysSuccess
=== RUN   TestLog_GetChildType/AlwaysFailure
=== RUN   TestLog_GetChildType/WhileSuccess
=== RUN   TestLog_GetChildType/WhileFailure
=== RUN   TestLog_GetChildType/Log
--- PASS: TestLog_GetChildType (0.00s)
    --- PASS: TestLog_GetChildType/nil_child (0.00s)
    --- PASS: TestLog_GetChildType/Action (0.00s)
    --- PASS: TestLog_GetChildType/Condition (0.00s)
    --- PASS: TestLog_GetChildType/Sequence (0.00s)
    --- PASS: TestLog_GetChildType/Selector (0.00s)
    --- PASS: TestLog_GetChildType/Parallel (0.00s)
    --- PASS: TestLog_GetChildType/Composite (0.00s)
    --- PASS: TestLog_GetChildType/Retry (0.00s)
    --- PASS: TestLog_GetChildType/Repeat (0.00s)
    --- PASS: TestLog_GetChildType/RepeatN (0.00s)
    --- PASS: TestLog_GetChildType/Invert (0.00s)
    --- PASS: TestLog_GetChildType/AlwaysSuccess (0.00s)
    --- PASS: TestLog_GetChildType/AlwaysFailure (0.00s)
    --- PASS: TestLog_GetChildType/WhileSuccess (0.00s)
    --- PASS: TestLog_GetChildType/WhileFailure (0.00s)
    --- PASS: TestLog_GetChildType/Log (0.00s)
=== RUN   TestLog_LoggingBehavior
=== RUN   TestLog_LoggingBehavior/Success_logging
2025/09/04 16:38:50 INFO Test logging child_status=Success child_type=Action
=== RUN   TestLog_LoggingBehavior/Failure_logging
2025/09/04 16:38:50 WARN Test logging child_status=Failure child_type=Action
=== RUN   TestLog_LoggingBehavior/Running_logging
=== RUN   TestLog_LoggingBehavior/Ready_logging
--- PASS: TestLog_LoggingBehavior (0.00s)
    --- PASS: TestLog_LoggingBehavior/Success_logging (0.00s)
    --- PASS: TestLog_LoggingBehavior/Failure_logging (0.00s)
    --- PASS: TestLog_LoggingBehavior/Running_logging (0.00s)
    --- PASS: TestLog_LoggingBehavior/Ready_logging (0.00s)
=== RUN   TestLog_CustomLogLevel
=== RUN   TestLog_CustomLogLevel/custom_error_level_for_success
2025/09/04 16:38:50 ERROR Custom level test child_status=Success child_type=Action
=== RUN   TestLog_CustomLogLevel/custom_info_level_for_failure
2025/09/04 16:38:50 INFO Info level failure child_status=Failure child_type=Action
=== RUN   TestLog_CustomLogLevel/custom_debug_level_for_running
=== RUN   TestLog_CustomLogLevel/nil_log_level_uses_defaults
2025/09/04 16:38:50 INFO Default level test child_status=Success child_type=Action
--- PASS: TestLog_CustomLogLevel (0.00s)
    --- PASS: TestLog_CustomLogLevel/custom_error_level_for_success (0.00s)
    --- PASS: TestLog_CustomLogLevel/custom_info_level_for_failure (0.00s)
    --- PASS: TestLog_CustomLogLevel/custom_debug_level_for_running (0.00s)
    --- PASS: TestLog_CustomLogLevel/nil_log_level_uses_defaults (0.00s)
=== RUN   TestLog_CustomLogLevelNoChild
2025/09/04 16:38:50 ERROR Log node has no child status=Failure
--- PASS: TestLog_CustomLogLevelNoChild (0.00s)
=== RUN   TestLog_StringWithLogLevel
=== RUN   TestLog_StringWithLogLevel/with_custom_error_level
2025/09/04 16:38:50 ERROR Error level test child_status=Success child_type=Action
=== RUN   TestLog_StringWithLogLevel/with_custom_debug_level
=== RUN   TestLog_StringWithLogLevel/no_custom_level_specified
2025/09/04 16:38:50 INFO Default level test child_status=Success child_type=Action
--- PASS: TestLog_StringWithLogLevel (0.00s)
    --- PASS: TestLog_StringWithLogLevel/with_custom_error_level (0.00s)
    --- PASS: TestLog_StringWithLogLevel/with_custom_debug_level (0.00s)
    --- PASS: TestLog_StringWithLogLevel/no_custom_level_specified (0.00s)
=== RUN   TestLog_ResetWithCustomLevel
2025/09/04 16:38:50 INFO Reset test with custom level child_status=Success child_type=Action
2025/09/04 16:38:50 INFO Log node reset message="Reset test with custom level"
--- PASS: TestLog_ResetWithCustomLevel (0.00s)
=== RUN   TestComplexBehaviorTreeWithLog
2025/09/04 16:38:50 WARN Monitoring critical action child_status=Failure child_type=Action
2025/09/04 16:38:50 WARN Monitoring critical action child_status=Failure child_type=Action
2025/09/04 16:38:50 INFO Monitoring critical action child_status=Success child_type=Action
--- PASS: TestComplexBehaviorTreeWithLog (0.00s)
PASS
ok  	github.com/rbrabson/behave	0.189s
