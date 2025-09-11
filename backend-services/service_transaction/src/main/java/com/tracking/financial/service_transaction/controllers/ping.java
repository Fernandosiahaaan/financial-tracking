package com.tracking.financial.service_transaction.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestParam;

@RestController
@RequestMapping
public class ping {
    @GetMapping("/ping")
    public String ping() {
        return new String("pong");
    }

}
