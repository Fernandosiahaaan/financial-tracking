package com.financial.tracking.service_category.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BaseResponse<T> {
    private boolean isSuccess;
    private String message;
    private T data;

    public static <T> BaseResponse<T> setResponse(boolean isSuccess, String message, T data) {
        return new BaseResponse<>(isSuccess, message, data);
    }
}
