package com.tracking.financial.service_transaction.repository;

import org.springframework.stereotype.Repository;

import com.tracking.financial.service_transaction.models.TransactionModels;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface TransactionRepository extends JpaRepository<TransactionModels, Long> {

} 
