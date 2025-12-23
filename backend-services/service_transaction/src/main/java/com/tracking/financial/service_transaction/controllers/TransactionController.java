package com.tracking.financial.service_transaction.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.tracking.financial.service_transaction.dto.BaseResponse;
import com.tracking.financial.service_transaction.dto.OnCreateRequest;
import com.tracking.financial.service_transaction.dto.TransactionCreateRequest;
import com.tracking.financial.service_transaction.dto.TransactionUpdateRequest;
import com.tracking.financial.service_transaction.services.TransactionService;

import jakarta.validation.Valid;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;

@RestController
@RequestMapping("/trns")
public class TransactionController {

    @Autowired
    private TransactionService transactionService;

    @PostMapping
    public ResponseEntity<BaseResponse> create(
            @Validated(OnCreateRequest.class) @RequestBody TransactionCreateRequest request) {
        return ResponseEntity.status(HttpStatus.OK).body(transactionService.create(request));
    }

    @GetMapping
    public ResponseEntity<BaseResponse> findAll() {
        return ResponseEntity.status(HttpStatus.OK).body(transactionService.findAll());
    }

    @GetMapping("/{id}")
    public ResponseEntity<BaseResponse> findByid(@Valid @PathVariable Long id) {
        return ResponseEntity.status(HttpStatus.OK).body(transactionService.findById(id));
    }

    @PutMapping
    public ResponseEntity<BaseResponse> update(@Valid @RequestBody TransactionUpdateRequest request) {
        return ResponseEntity.status(HttpStatus.OK).body(transactionService.update(request));
    }

    @DeleteMapping
    public ResponseEntity<BaseResponse> delete(@PathVariable Long id) {
        return ResponseEntity.status(HttpStatus.OK).body(transactionService.delete(id));
    }

}
