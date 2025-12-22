package com.tracking.financial.service_transaction.repository;

import org.springframework.stereotype.Repository;
import com.tracking.financial.service_transaction.models.TransactionModels;
import java.util.Optional;
import org.springframework.data.jpa.repository.JpaRepository;

@Repository
public interface TransactionRepository extends JpaRepository<TransactionModels, Long> {
    Optional<TransactionModels> findByName(String name);
}
