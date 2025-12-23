package com.tracking.financial.service_transaction.services;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionCreateRequest;
import com.tracking.financial.service_transaction.dto.TransactionUpdateRequest;

public interface TransactionService {
    BaseResponse create(TransactionCreateRequest item);

    BaseResponse findAll();

    BaseResponse findById(Long id);

    BaseResponse update(TransactionUpdateRequest item);

    BaseResponse delete(Long id);
}
