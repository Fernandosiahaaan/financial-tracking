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
    private TransactionRepository transactionRepository;

    @Override
    public BaseResponse create(TransactionRequest item) {
        List<TransactionModels> datas = transactionRepository.findAll();
        return BaseResponse.setResponse(true, "success", null, datas);
    }

    @Override
    public BaseResponse findAll() {
        return BaseResponse.setResponse(false, null, null, null);
    }

    @Override
    public BaseResponse findById(Long id) {
        return BaseResponse.setResponse(false, null, null, null);
    }

    @Override
    public BaseResponse update(TransactionRequest item) {
        return BaseResponse.setResponse(false, null, null, null);
    }

    @Override
    public BaseResponse delete(Long id) {
        return BaseResponse.setResponse(false, null, null, null);
    }

}
