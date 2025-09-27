namespace go billing_engine.data.loan_schedule_data
namespace java billing_engine.data.loan_schedule_data
namespace py billing_engine.data.loan_schedule_data

struct LoanScheduleSummaryData {
    1: optional i64 schedule_id,
    2: optional i32 week_number,
    3: optional i64 due_date,
    4: optional i64 amount,
    5: optional bool is_paid
}