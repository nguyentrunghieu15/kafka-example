# Schedule

- Ngày học lí thuyết : 16/03/2024 -> 20/03/2024
- Ngày viết code example : 20/03/2024 -> 21/03/2024 
- [Báo cáo](https://docs.google.com/document/d/1DMc82Tv70BuKPj2u6NrOeR4HU2S-bnOuAkDZWQ08_c8)

# Structure 
```
.
├── go.mod
├── go.sum
├── README.md
├── taxi-application
│   ├── estimate_cost_stream_service.go
│   ├── main.go
│   └── order_service.go
└── wikimedia_producer
    └── main.go

```
- taxi-application: triển khai 1 ứng dụng cơ bản thuê xe taxi trong đó khách hàng sẽ là 1 producer gửi yêu cầu order vaò topic order lên khụm kafka. 1 comsumer group sẽ lắng nghe topic order và đóng vai trò là người ước lượng giá thành , vì ước lượng giá thành là 1 công việc nặng bao gồm tính toán đường đi ngắn nhất, tính toán giá tiền, độ tắc đường, sau khi tính toán xong consumer sẽ push giá cả lên topic cost
- wikimedia_producer : triển khai 1 ứng dụng lắng nghe luồng stream event từ wikimedia, filter event và push lên kafka

