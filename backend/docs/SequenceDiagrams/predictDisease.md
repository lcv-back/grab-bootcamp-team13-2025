```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant MLService

    User->>Frontend: Nhập triệu chứng, bấm "Dự đoán bệnh"
    Frontend->>Backend: POST /api/symptoms/predict\n{symptoms, answers?, images?}
    Backend->>Backend: Xác thực JWT, lấy user_id
    alt Thiếu hoặc sai JWT
        Backend-->>Frontend: 401 Unauthorized (User ID not found)
        Frontend-->>User: Hiển thị lỗi đăng nhập
    else JWT hợp lệ
        Backend->>Backend: Parse và validate dữ liệu đầu vào
        alt Dữ liệu không hợp lệ
            Backend-->>Frontend: 400 Bad Request (invalid request)
            Frontend-->>User: Hiển thị lỗi dữ liệu
        else Dữ liệu hợp lệ
            Backend->>MLService: Gửi dữ liệu triệu chứng, ảnh, answers (nếu có)
            alt MLService lỗi
                MLService-->>Backend: Lỗi hoặc không phản hồi
                Backend-->>Frontend: 500 Internal Server Error (ML service error)
                Frontend-->>User: Hiển thị lỗi hệ thống
            else MLService trả về kết quả
                MLService-->>Backend: Danh sách bệnh, câu hỏi follow-up
                Backend->>Backend: Chuẩn hóa danh sách bệnh dự đoán (tên, xác suất, mô tả)
                Backend->>Backend: Kiểm tra có follow-up questions không?
                alt Có câu hỏi follow-up
                    Backend-->>Frontend: 200 OK, trả predicted_diseases, follow_up_questions
                    Frontend-->>User: Hiển thị danh sách bệnh dự đoán và các câu hỏi follow-up để user trả lời tiếp
                else Không có câu hỏi follow-up
                    Backend-->>Frontend: 200 OK, trả predicted_diseases, diagnosis completed
                    Frontend-->>User: Hiển thị kết quả chẩn đoán cuối cùng (tên bệnh, xác suất, mô tả)
                end
            end
        end
    end
```

---

**Chú thích:**
- Sơ đồ mô tả luồng dự đoán bệnh dựa trên triệu chứng, có thể gồm ảnh và câu trả lời follow-up.
- Thể hiện rõ các bước xác thực, xử lý lỗi, gọi ML service và trả kết quả.
- Áp dụng cho endpoint /api/symptoms/predict hoặc các endpoint tương tự trong hệ thống.
