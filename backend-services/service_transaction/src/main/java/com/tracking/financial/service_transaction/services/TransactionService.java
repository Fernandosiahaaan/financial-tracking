package com.tracking.financial.service_transaction.services;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionRequest;

public interface TransactionService {
    BaseResponse create(TransactionRequest item);

    BaseResponse findAll();

    BaseResponse findById(Long id);

    BaseResponse update(TransactionRequest item);

    BaseResponse delete(Long id);
}
