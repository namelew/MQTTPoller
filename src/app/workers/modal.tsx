'use client'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, Button, HStack, Input, FormControl, FormLabel, Checkbox, Flex } from "@chakra-ui/react";
import { useEffect, useState } from 'react';
import { IRequest } from '../interfaces/IExperiment';

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
        <Modal isOpen={openModal} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Parâmetros do Experimento</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <form>
                        <FormControl>
                            <FormLabel>Endereço do Broker</FormLabel>
                            <Input type="text" name="broker" value={formValues?.description.broker} onChange={handleChange} />
                            <FormLabel>Porta</FormLabel>
                            <Input type='number' name="port" value={formValues?.description.port} onChange={handleChange} />
                            <FormLabel>Versão do Protocolo</FormLabel>
                            <Input 
                                type='number'
                                name="mqttVersion"
                                value={formValues?.description.mqttVersion}
                                onChange={handleChange}
                            />
                            <FormLabel>Tópico</FormLabel>
                            <Input 
                                type='text'
                                name="topic"
                                value={formValues?.description.topic}
                                onChange={handleChange}
                            />
                            <FormLabel>Endereço Servidor NTP</FormLabel>
                            <Input 
                                type='text'
                                name="ntp"
                                value={formValues?.description.ntp}
                                onChange={handleChange}
                            />
                            <FormLabel>Usuário</FormLabel>
                            <Input 
                                type='text'
                                name="username"
                                value={formValues?.description.username}
                                onChange={handleChange}
                            />
                            <FormLabel>Senha</FormLabel>
                            <Input 
                                type='text'
                                name="password"
                                value={formValues?.description.password}
                                onChange={handleChange}
                            />
                            <FormLabel>Número de Mensagens Publicadas (Por Cliente) </FormLabel>
                            <Input 
                                type='number'
                                name="numMessages"
                                value={formValues?.description.numMessages}
                                onChange={handleChange}
                            />
                            <FormLabel>Tamanho das Mensagens (Em Bytes) </FormLabel>
                            <Input 
                                type='number'
                                name="payload"
                                value={formValues?.description.payload}
                                onChange={handleChange}
                            />
                            <FormLabel>Intervalo entra as Mensagens (Em Segundos) </FormLabel>
                            <Input 
                                type='number'
                                name="interval"
                                value={formValues?.description.interval}
                                onChange={handleChange}
                            />
                            <FormLabel>Número de Assinantes</FormLabel>
                            <Input 
                                type='number'
                                name="numSubscribers"
                                value={formValues?.description.numSubscribers}
                                onChange={handleChange}
                            />
                            <FormLabel>Número de Publicadores</FormLabel>
                            <Input 
                                type='number'
                                name="numPublishers"
                                value={formValues?.description.numPublishers}
                                onChange={handleChange}
                            />
                            <FormLabel>QoS das Assinaturas</FormLabel>
                            <Input 
                                type='number'
                                name="qosSubscriber"
                                value={formValues?.description.qosSubscriber}
                                onChange={handleChange}
                            />
                            <FormLabel>QoS das Publicações</FormLabel>
                            <Input 
                                type='number'
                                name="qosPublisher"
                                value={formValues?.description.qosPublisher}
                                onChange={handleChange}
                            />
                            <FormLabel>Tempo de Execução (Segundos)</FormLabel>
                            <Input 
                                type='number'
                                name="execTime"
                                value={formValues?.description.execTime}
                                onChange={handleChange}
                            />
                            <FormLabel>Timeout de Assinatura (Segundos)</FormLabel>
                            <Input 
                                type='number'
                                name="subscriberTimeout"
                                value={formValues?.description.subscriberTimeout}
                                onChange={handleChange}
                            />
                            <FormLabel>Tempo de Partida</FormLabel>
                            <Input 
                                type='number'
                                name="ramUp"
                                value={formValues?.description.ramUp}
                                onChange={handleChange}
                            />
                            <FormLabel>Tempo de Finalização</FormLabel>
                            <Input 
                                type='number'
                                name="rampDown"
                                value={formValues?.description.rampDown}
                                onChange={handleChange}
                            />
                            <Flex>
                                <FormLabel>Utilizar Assinatura Compatilhada?</FormLabel>
                                <Checkbox
                                    name="sharedSubscription"
                                    isChecked={formValues?.description.sharedSubscription}
                                    onChange={handleChange}
                                />
                            </Flex>
                            <Flex>
                                <FormLabel>Retenção?</FormLabel>
                                <Checkbox
                                    name="sharedSubscription"
                                    isChecked={formValues?.description.sharedSubscription}
                                    onChange={handleChange}
                                />
                            </Flex>
                            <Flex>
                                <FormLabel>Gerar arquivo de log?</FormLabel>
                                <Checkbox
                                    name="output"
                                    isChecked={formValues?.description.output}
                                    onChange={handleChange}
                                />
                            </Flex>
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