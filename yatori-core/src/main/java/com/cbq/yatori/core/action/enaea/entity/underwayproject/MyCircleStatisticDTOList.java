package com.cbq.yatori.core.action.enaea.entity.underwayproject;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@lombok.Data
public class MyCircleStatisticDTOList {
    @JsonProperty("list")
    private List<String> list;
}
