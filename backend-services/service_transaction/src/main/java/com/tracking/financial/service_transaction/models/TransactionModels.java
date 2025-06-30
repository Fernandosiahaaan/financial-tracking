package com.tracking.financial.service_transaction.models;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.UUID;

import org.hibernate.annotations.CreationTimestamp;

import lombok.Data;
import jakarta.persistence.Column;
import jakarta.persistence.Id;

public class TransactionModels {

    @Id
    private long id;

    @Column(length = 50, unique = true, nullable = false)
    private String name;

    @Column(name = "user_id", nullable = false)
    private UUID userId;
    
     @Column(name = "category_id", nullable = false)
    private UUID categoryId;

    private BigDecimal amount;
    
    private String description;

     @Column(name = "transaction_date", nullable = false)
    private LocalDate transactionDate;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, columnDefinition = "TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
    private LocalDateTime createdAt;

    @CreationTimestamp
    @Column(name = "updated_at", nullable = false, columnDefinition = "TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
    private LocalDateTime updadateAt;
}
