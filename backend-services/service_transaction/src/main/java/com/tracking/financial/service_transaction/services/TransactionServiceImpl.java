package com.tracking.financial.service_transaction.services;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.tracking.financial.service_transaction.clients.ClientCategory;
import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionRequest;
import com.tracking.financial.service_transaction.exceptions.BadRequestException;
import com.tracking.financial.service_transaction.exceptions.NotFoundException;
import com.tracking.financial.service_transaction.models.TransactionModels;
import com.tracking.financial.service_transaction.repository.TransactionRepository;

@Service
public class TransactionServiceImpl implements TransactionService {
    @Autowired
    private TransactionRepository transactionRepository;
    @Autowired
    private ClientCategory clientCategory;

    @Override
    public BaseResponse create(TransactionRequest item) {
        if (transactionRepository.findByName(item.getName()).isPresent()) {
            throw new BadRequestException("Transaction name '" + item.getName() + "' sudah tersedia");
        }

        if (!clientCategory.isExistCategoryById(item.getCategoryId().toString())) {
            throw new BadRequestException("Category tidak dikenal'");
        }

        TransactionModels data = new TransactionModels();
        data.setName(item.getName());
        data.setAmount(item.getAmount());
        data.setUserId(item.getUserId());
        data.setCategoryId(item.getCategoryId());
        data.setCreatedAt(LocalDateTime.now());
        TransactionModels output = transactionRepository.save(data);

        return BaseResponse.setResponse(true, "success", null, output);
    }

    @Override
    public BaseResponse findAll() {
        List<TransactionModels> datas = transactionRepository.findAll();
        return BaseResponse.setResponse(true, "success", null, datas);
    }

    @Override
    public BaseResponse findById(Long id) {
        TransactionModels data = transactionRepository.findById(id)
                .orElseThrow(() -> new NotFoundException("Category dengan id '" + id + "' tidak ditemukan"));
        ;
        return BaseResponse.setResponse(false, null, null, data);
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
