# Dùng redis lưu ID tối đa hiện tại
* URL gen ra ngắn (ID chạy từ 0)
* Problem: Sau 1 thời gian phải chạy lại routine set lại ID tối đa (Một số record sẽ bị xóa khi hết hạn, chạy để giảm số ID lại --> Link ngắn hơn)

