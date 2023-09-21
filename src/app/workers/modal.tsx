'use client'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, Button, HStack, FormControl } from "@chakra-ui/react";
import { useEffect, useState } from 'react';
import { mqttVersions, QoS } from "./static";
import { IRequest, IResult } from '../interfaces/IExperiment';
import TextArea from "./inputs/textarea";
import Number from "./inputs/number";
import Checkbox from "./inputs/checkbox";
import DropDown from "./inputs/dropdown";
import { api } from "../consumers/client";

interface Props {
    openModal: boolean,
    onClose: () => void,
    selected?: string[],
}

const ExperimentModal = ({ openModal, onClose, selected }: Props) => {
    const [formValues, setFormValues] = useState<IRequest>(
        {
            description: {
                tool: "mqtt-loader",
                broker: "",
                port:  0,
                mqttVersion:  3,
                numPublishers:  0,
                numSubscribers:  0,
                qosPublisher:  0,
                qosSubscriber:  0,
                sharedSubscription:  false,
                retain:  false,
                topic:  "",
                payload:  0,
                numMessages:  0,
                ramUp:  0,
                rampDown:  0,
                interval:  0,
                subscriberTimeout: 0,
                execTime: 0,
                logLevel: "INFO",
                ntp: "",
                output: true,
                username: "",
                password: "",
                tlsTruststore: "",
                tlsTruststorePass: "",
                tlsKeystore: "",
                tlsKeystorePass: ""
            }
        }
    );

    useEffect(() => {
        if (selected) {
            setFormValues((prevValues) => ({
                id: selected,
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

    const handleChange = (event:React.FormEvent<HTMLInputElement>) => {
        const { name, value } = event.currentTarget;

        setFormValues((prevValues) => ( prevValues ?
            {
                ...prevValues,
                description: { ...prevValues.description, [name]: value }
            } : prevValues));
    };

    const startExperiment = () => {
        api.post<IResult>("/experiment/start", formValues)
            .then((response) => {
                console.log(response.data);
            })
            .catch((error) => alert(error));
        onClose();
    }

    return (
        <Modal isOpen={openModal} onClose={onClose} size='3xl'>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Parâmetros do Experimento</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <form>
                        <FormControl>
                            <HStack justifyContent="space-between" width="100%">
                                <TextArea
                                    label="Endereço Broker"
                                    name="broker"
                                    value={formValues?.description.broker}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="Porta"
                                    name="port"
                                    value={formValues?.description.port}
                                    onChange={handleChange}
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
                                    onChange={handleChange}
                                />
                                <TextArea
                                    label="Usuário"
                                    name="username"
                                    value={formValues?.description.username}
                                    onChange={handleChange}
                                />
                                <TextArea
                                    label="Senha"
                                    name="password"
                                    value={formValues?.description.password}
                                    onChange={handleChange}
                                />
                            </HStack>
                            <HStack justifyContent="space-between" width="100%">
                                <Number
                                    label="N. Mensagens Publicadas"
                                    name="numMessages"
                                    value={formValues?.description.numMessages}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="Tamanho das Mensagens"
                                    name="payload"
                                    value={formValues?.description.payload}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="Intervalo entra Mensagens"
                                    name="interval"
                                    value={formValues?.description.interval}
                                    onChange={handleChange}
                                />
                            </HStack>
                            <HStack justifyContent="space-between" width="100%">
                                <TextArea label="Tópico" name="topic" value={formValues?.description.topic} onChange={handleChange}/>
                                <Number
                                    label="N. Publicadores"
                                    name="numPublishers"
                                    value={formValues?.description.numPublishers}
                                    onChange={handleChange}
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
                                    label="N. Assinantes"
                                    name="numSubscribers"
                                    value={formValues?.description.numSubscribers}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="Timeout Assinatura"
                                    name="subscriberTimeout"
                                    value={formValues?.description.subscriberTimeout}
                                    onChange={handleChange}
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
                                    onChange={handleChange}
                                />
                                <Number 
                                    label='Tempo de Partida'
                                    name="ramUp"
                                    value={formValues?.description.ramUp}
                                    onChange={handleChange}
                                />
                                <Number 
                                    label='Tempo de Finalização'
                                    name="rampDown"
                                    value={formValues?.description.rampDown}
                                    onChange={handleChange}
                                />
                            </HStack>
                            <Checkbox
                                label="Utilizar Assinatura Compatilhada?"
                                name="sharedSubscription"
                                isChecked={formValues?.description.sharedSubscription}
                                disabled={formValues?.description.mqttVersion !== 5}
                                onChange={handleChange}
                            />
                            <Checkbox
                                label="Retenção?"
                                name="retain"
                                isChecked={formValues?.description.retain}
                                onChange={handleChange}
                            />
                            <Checkbox
                                label="Gerar arquivo de log?"
                                name="output"
                                isChecked={formValues?.description.output}
                                onChange={handleChange}
                            />
                        </FormControl>
                    </form>
                </ModalBody>
                <ModalFooter>
                    <HStack justifyContent="space-between" width="100%">
                        <Button colorScheme="blue" mr={3} onClick={onClose}>
                            Fechar
                        </Button>
                        <Button colorScheme="green" mr={3} onClick={startExperiment}>
                            Iniciar
                        </Button>
                    </HStack>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default ExperimentModal;