package com.tracking.financial.service_transaction.services;

import java.util.List;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionRequest;
import com.tracking.financial.service_transaction.models.TransactionModels;
import com.tracking.financial.service_transaction.repository.TransactionRepository;

@Service
public class TransactionServiceImpl implements TransactionService {
    @Autowired
    private TransactionRepository repository;

    @Override
    public BaseResponse<TransactionModels> create(TransactionRequest item) {
        return BaseResponse.setResponse(false, null, null);
    }

    @Override
    public BaseResponse<List<TransactionModels>> findAll() {
        return BaseResponse.setResponse(false, null, null);
    }

    @Override
    public BaseResponse<TransactionModels> findById(Long id) {
        return BaseResponse.setResponse(false, null, null);
    }

    @Override
    public BaseResponse<TransactionModels> update(TransactionRequest item) {
        return BaseResponse.setResponse(false, null, null);
    }

    @Override
    public BaseResponse<Void> delete(Long id) {
        return BaseResponse.setResponse(false, null, null);
    }

}
