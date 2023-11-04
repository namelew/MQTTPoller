'use client'
import { Button, HStack, FormControl } from "@chakra-ui/react";
import { useEffect, useState } from 'react';
import { mqttVersions, QoS } from "../static";
import { IRequest } from 'interfaces/IExperiment';
import TextArea from "./inputs/textarea";
import Number from "./inputs/number";
import Checkbox from "./inputs/checkbox";
import DropDown from "./inputs/dropdown";
import { useMutation } from "@tanstack/react-query";
import { startExperiment } from "consumer";
import Modal from "components/modal";

interface Props {
    openModal: boolean,
    onClose: () => void,
    selected?: string[],
}

const ExperimentModal = ({ openModal, onClose, selected }: Props) => {
    const defaultParams = {
        description: {
            tool: "mqtt-loader",
            broker: "localhost",
            port:  1883,
            mqttVersion:  3,
            numPublishers:  0,
            numSubscribers:  0,
            qosPublisher:  0,
            qosSubscriber:  0,
            sharedSubscription:  false,
            retain:  false,
            topic:  "mqtt-test-topic",
            payload:  0,
            numMessages:  0,
            ramUp:  0,
            rampDown:  0,
            interval:  0,
            subscriberTimeout: 0,
            execTime: 5,
            logLevel: "INFO",
            ntp: "a.st1.ntp.br",
            username: "",
            password: "",
            tlsTruststore: "",
            tlsTruststorePass: "",
            tlsKeystore: "",
            tlsKeystorePass: ""
        }
    }

    const submitMutation = useMutation(startExperiment, {
        onSuccess: (result) => console.log(result),
        onError: (error) => console.log(error)
    });

    const [formValues, setFormValues] = useState<IRequest>(defaultParams);

    useEffect(() => {
        if (selected) {
            setFormValues((prevValues) => ({
                id: selected,
                attempts: 0,
                description: {
                    ...prevValues?.description,
                    tool: "mqttloader"
                }
            }));
        }
    }, [selected]);

    const handlerSelectChange = (event: React.FormEvent<HTMLSelectElement>) => {
        const { name, value } = event.currentTarget;
        setFormValues((prevValues) => ( prevValues ?
            {
                ...prevValues,
                description: { ...prevValues.description, [name]: +value }
            } : prevValues));
    }

    const handleChangeNumber = (event:React.FormEvent<HTMLInputElement>) => {
        const { name, value } = event.currentTarget;

        setFormValues((prevValues) => ( prevValues ?
            {
                ...prevValues,
                description: { ...prevValues.description, [name]: +value }
            } : prevValues));
    };

    const handleChangeBool = (event:React.FormEvent<HTMLInputElement>) => {
        const { name } = event.currentTarget;

        switch (name) {
            case "sharedSubscription":
                setFormValues((prevValues) => ( prevValues ?
                    {
                        ...prevValues,
                        description: { ...prevValues.description, sharedSubscription: !prevValues.description.sharedSubscription }
                    } : prevValues));
                break;
            case "retain":
                setFormValues((prevValues) => ( prevValues ?
                    {
                        ...prevValues,
                        description: { ...prevValues.description, retain: !prevValues.description.retain }
                    } : prevValues));
                break;
        }
    };

    const handleChangeString = (event:React.FormEvent<HTMLInputElement>) => {
        const { name, value } = event.currentTarget;

        setFormValues((prevValues) => ( prevValues ?
            {
                ...prevValues,
                description: { ...prevValues.description, [name]: value }
            } : prevValues));
    };

    const onSubmit = () => {
        submitMutation.mutate(formValues, {
            onSuccess: (response) => {
                console.log(response);
            },
            onError: (error, variables, context) => {
                console.log(error, variables, context);
            }
        });
        setFormValues({
            id: selected,
            attempts: 0,
            description: defaultParams.description,
        });
        onClose();
    }

    return (
        <Modal title="Parâmetros do Experimento" isOpen={openModal} onClose={onClose} size='3xl' footer={
            <HStack justifyContent="space-between" width="100%">
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Fechar
                </Button>
                <Button colorScheme="green" mr={3} onClick={onSubmit}>
                    Iniciar
                </Button>
            </HStack>
        }>
            <form>
                <FormControl>
                    <HStack justifyContent="space-between" width="100%">
                        <TextArea
                            label="Endereço Broker"
                            name="broker"
                            value={formValues?.description.broker}
                            onChange={handleChangeString}
                        />
                        <Number
                            label="Porta"
                            name="port"
                            value={formValues?.description.port}
                            onChange={handleChangeNumber}
                        />
                        <DropDown
                            label="Versão do Protocolo"
                            name="mqttVersion"
                            value={formValues?.description.mqttVersion}
                            onChange={handlerSelectChange}
                            options={mqttVersions}
                        />
                    </HStack>
                    <HStack justifyContent="space-between" width="100%">
                        <TextArea
                            label="Endereço Servidor NTP"
                            name="ntp"
                            value={formValues?.description.ntp}
                            onChange={handleChangeString}
                        />
                        <TextArea
                            label="Usuário"
                            name="username"
                            value={formValues?.description.username}
                            onChange={handleChangeString}
                        />
                        <TextArea
                            label="Senha"
                            name="password"
                            value={formValues?.description.password}
                            onChange={handleChangeString}
                        />
                    </HStack>
                    <HStack justifyContent="space-between" width="100%">
                        <Number
                            label="N. Mensagens Publicadas"
                            name="numMessages"
                            value={formValues?.description.numMessages}
                            onChange={handleChangeNumber}
                        />
                        <Number
                            label="Tamanho das Mensagens"
                            name="payload"
                            value={formValues?.description.payload}
                            onChange={handleChangeNumber}
                        />
                        <Number
                            label="Intervalo entra Mensagens"
                            name="interval"
                            value={formValues?.description.interval}
                            onChange={handleChangeNumber}
                        />
                    </HStack>
                    <HStack justifyContent="space-between" width="100%">
                        <TextArea label="Tópico" name="topic" value={formValues?.description.topic} onChange={handleChangeString}/>
                        <Number
                            label="N. Publicadores"
                            name="numPublishers"
                            value={formValues?.description.numPublishers}
                            onChange={handleChangeNumber}
                        />
                        <DropDown
                            label="QoS Publicações"
                            name="qosPublisher"
                            value={formValues?.description.qosPublisher}
                            onChange={handlerSelectChange}
                            options={QoS}
                        />
                    </HStack>
                    <HStack justifyContent="space-between" width="100%">
                        <Number
                            label="Timeout Assinatura"
                            name="subscriberTimeout"
                            value={formValues?.description.subscriberTimeout}
                            onChange={handleChangeNumber}
                        />
                        <Number
                            label="N. Assinantes"
                            name="numSubscribers"
                            value={formValues?.description.numSubscribers}
                            onChange={handleChangeNumber}
                        />
                        <DropDown
                            label="QoS Assinaturas"
                            name="qosSubscriber"
                            value={formValues?.description.qosSubscriber}
                            onChange={handlerSelectChange}
                            options={QoS}
                        />
                    </HStack>
                    <HStack justifyContent="space-between" width="100%">
                        <Number 
                            label='Tempo de Execução'
                            name="execTime"
                            value={formValues?.description.execTime}
                            onChange={handleChangeNumber}
                        />
                        <Number 
                            label='Tempo de Partida'
                            name="ramUp"
                            value={formValues?.description.ramUp}
                            onChange={handleChangeNumber}
                        />
                        <Number 
                            label='Tempo de Finalização'
                            name="rampDown"
                            value={formValues?.description.rampDown}
                            onChange={handleChangeNumber}
                        />
                    </HStack>
                    <Checkbox
                        label="Utilizar Assinatura Compatilhada?"
                        name="sharedSubscription"
                        isChecked={formValues?.description.sharedSubscription}
                        disabled={formValues?.description.mqttVersion !== 5}
                        onChange={handleChangeBool}
                    />
                    <Checkbox
                        label="Retenção?"
                        name="retain"
                        isChecked={formValues?.description.retain}
                        onChange={handleChangeBool}
                    />
                </FormControl>
            </form>
        </Modal>
    );
};

export default ExperimentModal;