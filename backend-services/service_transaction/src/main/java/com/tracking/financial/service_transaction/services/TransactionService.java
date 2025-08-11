package com.tracking.financial.service_transaction.services;

import java.util.List;
import java.util.UUID;

import org.apache.naming.TransactionRef;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionRequest;
import com.tracking.financial.service_transaction.models.TransactionModels;

public interface TransactionService {
    BaseResponse<TransactionModels> create(TransactionRequest item);
    BaseResponse<List<TransactionModels>> findAll();
    BaseResponse<TransactionModels> findById(Long id);
    BaseResponse<TransactionModels> update(TransactionRequest item);
    BaseResponse<Void> delete(Long id);
}
