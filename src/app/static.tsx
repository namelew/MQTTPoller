export const mqttVersions = [
    {
      key: "v3.1",
      name: "MQTT v3.1",
      value: 3  
    },
    {
        key: "v5",
        name: "MQTT v5",
        value: 5
    }
]

export const QoS = [
    {
        key: "QoS 0",
        name: "Melhor Esforço (QoS 0)",
        value: 0
    },
    {
        key: "QoS 1",
        name: "Pelo menos um (QoS 1)",
        value: 1
    },
    {
        key: "QoS 2",
        name: "Exatamente um (QoS 2)",
        value: 2
    }
]

export const logLevels = [
    {
        key: "info",
        name: "Informações",
        value: "INFO"  
    },
    {
        key: "all",
        name: "Tudo",
        value: "ALL"  
    },
    {
        key: "warn",
        name: "Avisos",
        value: "WARNING"  
    },
    {
        key: "severe",
        name: "Apenas Erros",
        value: "SEVERE"  
    }
]