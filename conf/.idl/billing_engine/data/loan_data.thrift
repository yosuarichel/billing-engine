namespace go billing_engine.data.product_data
namespace java billing_engine.data.product_data
namespace py billing_engine.data.product_data

struct LoanData {
   1: optional string loan_id,
   2: optional string customer_id,
   3: optional i64 principal,
   4: optional double interest_rate,
   5: optional i64 total_amount,
   6: optional i32 term_weeks,
   7: optional i64 start_date,
   8: optional string status,
}