export interface IRequest {
    id: string[],
    attempts: number,
    description: {
        tool: string,
        broker: string,
        port:  number,
        mqttVersion:  number,
        numPublishers:  number,
        numSubscribers:  number,
        qosPublisher:  number,
        qosSubscriber:  number,
        sharedSubscription:  boolean,
        retain:  boolean,
        topic:  string,
        payload:  number,
        numMessages:  number,
        ramUp:  number,
        rampDown:  number,
        interval:  number,
        subscriberTimeout:  number,
        execTime:  number,
        logLevel:  string,
        ntp:  string,
        output:  boolean,
        username:  string,
        password:  string,
        tlsTruststore:  string,
        tlsTruststorePass:  string,
        tlsKeystore:  string,
        tlsKeystorePass:  string
    }
};

export interface IResult {
    meta: {
        id: number,
        error: string,
        tool: string,
        literal: string,
        log_file: {
            name: string,
            data: File,
            extension: string
        },
    },
    publish: {
        max_throughput: Number,
        avg_throughput: Number,
        publiqued_messages: number,
        per_second_throungput: number[]
    },
    subscribe: {
        max_throughput: Number,
        avg_throughput: Number,
        received_messages: number,
        per_second_throungput: number[],
        latency: Number,
        avg_latency: Number
    }
};

export interface IExperiment {
    id: number,
    broker: string,
    port:  number,
    mqttVersion:  number,
    numPublishers:  number,
    numSubscribers:  number,
    qosPublisher:  number,
    qosSubscriber:  number,
    sharedSubscription:  boolean,
    retain:  boolean,
    topic:  string,
    payload:  number,
    numMessages:  number,
    ramUp:  number,
    rampDown:  number,
    interval:  number,
    subscriberTimeout:  number,
    execTime:  number,
    logLevel:  string,
    ntp:  string,
    output:  boolean,
    username:  string,
    password:  string,
    tlsTruststore:  string,
    tlsTruststorePass:  string,
    tlsKeystore:  string,
    tlsKeystorePass:  string,
    finish:boolean,
    error: string,
    results: IResult[],
    workers: string[]
};