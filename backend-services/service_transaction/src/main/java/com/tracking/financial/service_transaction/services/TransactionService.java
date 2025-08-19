package com.tracking.financial.service_transaction.services;

import java.util.List;
import java.util.UUID;

import org.apache.naming.TransactionRef;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionRequest;
import com.tracking.financial.service_transaction.models.TransactionModels;

public interface TransactionService {
    BaseResponse create(TransactionRequest item);

    BaseResponse findAll();

    BaseResponse findById(Long id);

    BaseResponse update(TransactionRequest item);

    BaseResponse delete(Long id);
}
