package com.tracking.financial.service_transaction.dto;

import java.math.BigDecimal;
import java.util.UUID;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import jakarta.validation.constraints.DecimalMin;
import com.fasterxml.jackson.annotation.JsonProperty;

@AllArgsConstructor
@Builder
@Getter
public class TransactionRequest {

    @JsonProperty("id")
    private Long id;

    @JsonProperty("name")
    @NotBlank(message = "Name is required")
    private String name;

    @JsonProperty("userId")
    @NotNull(message = "user_id is required")
    private UUID userId;

    @JsonProperty("amount")
    @NotNull(message = "amount is required")
    @DecimalMin(value = "0.01", message = "amount must be greater than 0")
    private BigDecimal amount;

    @JsonProperty("categoryId")
    @NotNull(groups = OnCreateRequest.class, message = "category_id is required")
    private UUID categoryId;

    @JsonProperty("description")
    private String description;

    // Getter & Setter
}