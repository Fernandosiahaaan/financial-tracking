package com.tracking.financial.service_transaction.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BaseResponse<T> {
    private boolean isSuccess;
    private String message;
    private String messageError;
    private T data;

    public static <T> BaseResponse<T> setResponse(boolean isSuccess, String message, String messageError) {
        return new BaseResponse<>(isSuccess, message, messageError, null);
    }
}
