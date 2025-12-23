package com.tracking.financial.service_transaction.dto;

import java.math.BigDecimal;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;

@AllArgsConstructor
@Builder
@Getter
public class TransactionUpdateRequest {
    
    @JsonProperty("id")
    @NotNull(message = "Id is required")
    private Long id;

    @JsonProperty("name")
    private String name;

    @JsonProperty("userId")
    private UUID userId;

    @JsonProperty("amount")
    private BigDecimal amount;

    @JsonProperty("categoryId")
    private UUID categoryId;

    @JsonProperty("description")
    private String description;
}
