'use client'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, Button, HStack, FormControl, Stack } from "@chakra-ui/react";
import { useEffect, useState } from 'react';
import { IRequest } from '../interfaces/IExperiment';
import TextArea from "./inputs/textarea";
import Number from "./inputs/number";
import Checkbox from "./inputs/checkbox";

interface Props {
    openModal: boolean,
    onClose: () => void,
    selected?: string[],
}

const ExperimentModal = ({ openModal, onClose, selected }: Props) => {
    const [formValues, setFormValues] = useState<IRequest>();

    useEffect(() => {
        if (selected) {
            setFormValues( (prevValues) => ( prevValues ?
                {
                    id: selected,
                    attempts: 0,
                    description: {
                        ...prevValues.description, tool: "mqtt-loader"
                    }
                } : prevValues));
        }
    }, [selected]);

    const handleChange = (event:React.FormEvent<HTMLInputElement>) => {
        const { name, value } = event.currentTarget;
        setFormValues((prevValues) => ( prevValues ?
            {
                ...prevValues,
                description: { ...prevValues.description, [name]: value }
            } : prevValues));
    };

    return (
        <Modal isOpen={openModal} onClose={onClose} size='3xl'>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Parâmetros do Experimento</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <form>
                        <FormControl>
                            <Stack direction="row" align="center" justify="space-between">
                                <TextArea
                                    label="Endereço do Broker"
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
                                <Number
                                    label="Versão do Protocolo"
                                    name="mqttVersion"
                                    value={formValues?.description.mqttVersion}
                                    onChange={handleChange}
                                />
                            </Stack>
                            <Stack direction="row" align="center" justify="space-between">
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
                            </Stack>
                            <Stack direction="row" align="center" justify="space-between">
                                <Number
                                    label="Número de Mensagens Publicadas"
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
                                    label="Intervalo entra as Mensagens"
                                    name="interval"
                                    value={formValues?.description.interval}
                                    onChange={handleChange}
                                />
                            </Stack>
                            <Stack direction="row" align="center" justify="space-between">
                                <TextArea label="Tópico" name="topic" value={formValues?.description.topic} onChange={handleChange}/>
                                <Number
                                    label="N. Publicadores"
                                    name="numPublishers"
                                    value={formValues?.description.numPublishers}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="QoS Publicações"
                                    name="qosPublisher"
                                    value={formValues?.description.qosPublisher}
                                    onChange={handleChange}
                                />
                            </Stack>
                            <Stack direction="row" align="center" justify="space-between">
                                <Number
                                    label="N. Assinantes"
                                    name="numSubscribers"
                                    value={formValues?.description.numSubscribers}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="QoS Assinaturas"
                                    name="qosSubscriber"
                                    value={formValues?.description.qosSubscriber}
                                    onChange={handleChange}
                                />
                                <Number
                                    label="Timeout Assinatura"
                                    name="subscriberTimeout"
                                    value={formValues?.description.subscriberTimeout}
                                    onChange={handleChange}
                                />
                            </Stack>
                            <Stack direction="row" align="center" justify="space-between">
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
                            </Stack>
                            <Checkbox
                                label="Utilizar Assinatura Compatilhada?"
                                name="sharedSubscription"
                                isChecked={formValues?.description.sharedSubscription}
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
                        <Button colorScheme="green" mr={3} onClick={onClose}>
                            Iniciar
                        </Button>
                    </HStack>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default ExperimentModal;