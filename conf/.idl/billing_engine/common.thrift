namespace go billing_engine.common
namespace java billing_engine.common
namespace py billing_engine.common

// GET /api/items?pageSize=20&pageNum=3
struct PageNumberPagination {
   1: i32 page_num,
   2: i32 page_size,
   3: i32 total,
}