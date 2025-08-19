package com.tracking.financial.service_transaction.dto;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class BaseResponse {
    private boolean isSuccess;
    private String message;
    private String messageError;
    private Object data;

    public BaseResponse(boolean isSuccess, String message, String messageError, Object data) {
        this.isSuccess = isSuccess;
        this.message = message;
        this.messageError = messageError;
        this.data = data;
    }

    public static BaseResponse setResponse(boolean isSuccess, String message, String messageError, Object data) {
        return new BaseResponse(isSuccess, message, messageError, data);
    }
}
