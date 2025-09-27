namespace go billing_engine
namespace java billing_engine
namespace py billing_engine

include "../base.thrift"
include "./data/loan_schedule_data.thrift"
include "./data/billing_data.thrift"
include "./data/payment_data.thrift"
include "./common.thrift"

struct IsDelinquentRequest {
    1: optional string loan_id;
    

    100: optional i32 version (api.path = "version");
    255: optional base.Base Base;
}

struct IsDelinquentResponse {
    1: optional bool is_delinquent;

    255: base.BaseResp BaseResp;
}

struct CreateCustomerRequest {
    1: required string name;
    2: required string phone_number;
    

    100: optional i32 version (api.path = "version");
    255: optional base.Base Base;
}

struct CreateCustomerResponse {
    1: optional string customer_id;

    255: base.BaseResp BaseResp;
}

struct CreateLoanRequest {
    1: required string customer_id,
    2: required i64 principal,
    3: required i32 term_weeks,

    100: optional i32 version (api.path = "version");
    255: optional base.Base Base;
}

struct CreateLoanResponse {
    1: optional string loan_id,
    2: optional string customer_id,
    3: optional i64 principal,
    4: optional i64 total_amount,
    5: optional i32 term_weeks,
    6: optional i64 start_date,
    7: optional string status,
    8: optional list<loan_schedule_data.LoanScheduleSummaryData> schedules

    255: base.BaseResp BaseResp;
}

struct GetOutstandingRequest {
    1: optional string loan_id;
    

    100: optional i32 version (api.path = "version");
    255: optional base.Base Base;
}

struct GetOutstandingResponse {
    1: optional billing_data.OutstandingData data;

    255: base.BaseResp BaseResp;
}

struct MakePaymentRequest {
    1: optional string loan_id;
    2: optional i64 amount;
    

    100: optional i32 version (api.path = "version");
    255: optional base.Base Base;
}

struct MakePaymentResponse {
    1: optional payment_data.PaymentData data;

    255: base.BaseResp BaseResp;
}

// Service
service BillingEngineService {
    // Customer
    CreateCustomerResponse CreateCustomer(1: CreateCustomerRequest req) (api.post="/api/v1/customer/create", agw.preserve_base = 'true')

    // Billing
    IsDelinquentResponse IsDelinquent(1: IsDelinquentRequest req) (api.post="/api/v1/billing/is_delinquent", agw.preserve_base = 'true')
    CreateLoanResponse CreateLoan(1: CreateLoanRequest req) (api.post="/api/v1/billing/loan/create", agw.preserve_base = 'true')
    GetOutstandingResponse GetOutstanding(1: GetOutstandingRequest req) (api.post="/api/v1/billing/get_outstanding", agw.preserve_base = 'true')
    MakePaymentResponse MakePayment(1: MakePaymentRequest req) (api.post="/api/v1/billing/payment/create", agw.preserve_base = 'true')
}

