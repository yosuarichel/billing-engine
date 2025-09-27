namespace go billing_engine.data.payment_data
namespace java billing_engine.data.payment_data
namespace py billing_engine.data.payment_data

struct PaymentData {
   1: optional string payment_id,
   2: optional i64 amount,
   3: optional i64 pay_amount,
   4: optional i64 outstanding,
   5: optional i32 term_weeks,
   6: optional i32 weeks_remain,
}