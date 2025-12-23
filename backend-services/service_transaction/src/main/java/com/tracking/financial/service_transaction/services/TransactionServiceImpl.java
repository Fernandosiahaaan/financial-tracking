package com.tracking.financial.service_transaction.services;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.JsonSerializable.Base;
import com.tracking.financial.service_transaction.clients.ClientCategory;
import com.tracking.financial.service_transaction.clients.ClientUser;
import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.TransactionCreateRequest;
import com.tracking.financial.service_transaction.dto.TransactionUpdateRequest;
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
    @Autowired
    private ClientUser userCategory;

    @Override
    public BaseResponse create(TransactionCreateRequest item) {
        if (transactionRepository.findByName(item.getName()).isPresent()) {
            throw new BadRequestException("Transaction name '" + item.getName() + "' sudah tersedia");
        }

        // //TODO: set auth for service-user
        // if (!userCategory.isExistUserById(item.getCategoryId().toString())) {
        //     throw new BadRequestException("User tidak dikenal'");
        // }

        if (!clientCategory.isExistCategoryById(item.getCategoryId().toString())) {
            throw new BadRequestException("Category tidak dikenal'");
        }

        // Input Data to Insert
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
                .orElseThrow(() -> new NotFoundException("Transaction id '" + id + "' tidak ditemukan"));
        return BaseResponse.setResponse(false, null, null, data);
    }

    @Override
    public BaseResponse update(TransactionUpdateRequest item) {
        // Check Id Transaction
        TransactionModels existing = transactionRepository.findById(item.getId())
                .orElseThrow(() -> new NotFoundException("Transaction id '" + item.getId() + "' tidak ditemukan"));
        
        // Input Name
        if ( item.getName() != null && !item.getName().isEmpty() && !item.getName().isBlank()) {
            existing.setName(item.getName());
        }

        // Input Category
        if ( item.getCategoryId() != null && !item.getCategoryId().toString().isEmpty() && !item.getCategoryId().toString().isBlank()) {
            // validate category id
            if (!clientCategory.isExistCategoryById(item.getCategoryId().toString())) {
                throw new BadRequestException("Category tidak dikenal'");
            }

            existing.setCategoryId(item.getCategoryId());
        }

        // Input Amount
        if ( item.getAmount() != null && item.getAmount().intValue() > 0) {
            existing.setAmount(item.getAmount());
        }

        // Input Description
        if ( item.getDescription() != null && !item.getDescription().isEmpty() && !item.getDescription().isBlank()) {
            existing.setDescription(item.getDescription());
        }

        // Input Updated At
        existing.setUpdadateAt(LocalDateTime.now());
        
        TransactionModels data = transactionRepository.save(existing);
        return BaseResponse.setResponse(true, "Success", null, data);
    }

    @Override
    public BaseResponse delete(Long id) {
        return BaseResponse.setResponse(false, null, null, null);
    }

}
