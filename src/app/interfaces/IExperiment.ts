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
        tlsKeystorePass:  string,
    }
};

export interface IResult {
    meta: {
        id: any,
        error: any,
        tool: any,
        literal: any,
        log_file: {
            name: string,
            data: File,
            extension: string
        },
    },
    publish: {
        max_throughput: any,
        avg_throughput: any,
        publiqued_messages: any,
        per_second_throungput: any,
    },
    subscribe: {
        max_throughput: any,
        avg_throughput: any,
        received_messages: any,
        per_second_throungput: any,
        latency: any,
        avg_latency: any,
    },
};

export interface IExperiment {

};