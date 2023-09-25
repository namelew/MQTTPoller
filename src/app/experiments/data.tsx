'use client'
import { Tr, Td, Button, Flex, Text, HStack, Box, Container } from "@chakra-ui/react";
import { useMutation } from "@tanstack/react-query";
import Modal from "components/modal";
import { deleteExperiment } from "consumer";
import { IExperiment } from "interfaces/IExperiment";
import { useState } from "react";
import { QoS, mqttVersions } from "static";
import ResultContainer from "./results";
import Carousel from "components/carossel";

interface Props {
    experiments?:IExperiment[],
}

const Experiments = ( { experiments } : Props) => {
    const [selectedModal, setSelectedModal] = useState<number | null>(null);


    const onDelete = useMutation(deleteExperiment, {
        onSuccess: () => alert("Experimento excluido com sucesso"),
        onError: (error) => console.log(error)
    });

    return (
        <>
            {experiments?.map((experiment) => {
                const mqttVersion = mqttVersions.find((version) => version.value === experiment.mqttVersion);
                const mqttVersionName = mqttVersion ? mqttVersion.name : 'Unknown version';

                const qosPub = QoS.find((quality) => quality.value === experiment.qosPublisher);
                const qosPubName = qosPub ? qosPub.name : 'Unknown QoS';

                const qosSub = QoS.find((quality) => quality.value === experiment.qosSubscriber);
                const qosSubName = qosSub ? qosSub.name : 'Unknown QoS';

                return (
                    <Tr key={experiment.id}>
                        <Td>{experiment.id}</Td>
                        <Td>{experiment.broker}:{experiment.port}</Td>
                        <Td>{mqttVersionName}</Td>
                        <Td>{experiment.topic}</Td>
                        <Td>{experiment.execTime} s</Td>
                        <Td>{experiment.finish ? 'Sim' : 'Não'}</Td>
                        <Td>{experiment.error !== '' ? 'Sim' : 'Não'}</Td>
                        <Td>
                            <Button colorScheme="teal" size="sm" variant="outline" onClick={() => setSelectedModal(experiment.id)}>
                                Detalhar
                            </Button>
                            <Modal
                                isOpen={selectedModal === experiment.id} 
                                onClose={() => setSelectedModal(null)}
                                title='Detalhes do Experimento'
                                size='5xl'
                            >
                                <Flex gap='1' p={2}>
                                    <Text fontWeight="bold">ID:</Text>
                                    <Text>{experiment.id}</Text>
                                </Flex>
                                <HStack>
                                    <Container>
                                        <Text fontWeight="bold" textAlign="center">Parâmetros</Text>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Broker:</Text>
                                            <Text>{experiment.broker}:{experiment.port}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">NTP:</Text>
                                            <Text>{experiment.ntp === "" ? 'Não informado' : experiment.ntp}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Versão do Protocolo:</Text>
                                            <Text>{mqttVersionName}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Tópico:</Text>
                                            <Text>{experiment.topic}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Usuário:</Text>
                                            <Text>{experiment.username === "" ? 'Não informado' : experiment.username}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Senha:</Text>
                                            <Text>{experiment.password === "" ? 'Não informado' : experiment.password}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Num. Publicadores:</Text>
                                            <Text>{experiment.numPublishers}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Num. de Assinantes:</Text>
                                            <Text>{experiment.numSubscribers}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">QoS de Publicação:</Text>
                                            <Text>{qosPubName}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">QoS de Assinatura:</Text>
                                            <Text>{qosSubName}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Timeout de Assinatura:</Text>
                                            <Text>{experiment.subscriberTimeout} s</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Num. Mensagens:</Text>
                                            <Text>{experiment.numMessages}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Tamanho Mensagens:</Text>
                                            <Text>{experiment.payload} bytes</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Intervalo entre as mensagens:</Text>
                                            <Text>{experiment.interval} s</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Tempo de Execução:</Text>
                                            <Text>{experiment.execTime} s</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Tempo de Partida:</Text>
                                            <Text>{experiment.ramUp} s</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Tempo de Finalização:</Text>
                                            <Text>{experiment.rampDown} s</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Utilizou Assinatura Compartilhada?</Text>
                                            <Text>{experiment.sharedSubscription ? 'Sim' : 'Não'}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Utilizou Marcador de Retenção?</Text>
                                            <Text>{experiment.retain ? 'Sim' : 'Não'}</Text>
                                        </Flex>
                                        <Flex gap='1'>
                                            <Text fontWeight="bold">Finalizado?</Text>
                                            <Text>{experiment.finish ? 'Sim' : 'Não'}</Text>
                                        </Flex>
                                        {experiment.error !== '' && 
                                            <Flex gap='1'>
                                                <Text fontWeight="bold">Error:</Text>
                                                <Text>{experiment.error}</Text>
                                            </Flex>
                                        }
                                    </Container>
                                    <Container>
                                        <Text fontWeight="bold" textAlign="center">Resultados</Text>
                                        {experiment.results === null || experiment.results.length === 0 ? 
                                            'Sem resultados' :
                                            <Carousel
                                                items={experiment.results.map( (result, index) => <ResultContainer key={index} result={result} workerID={experiment.workers[index]}/>)}
                                            />
                                        }
                                    </Container>
                                </HStack>
                            </Modal>
                        </Td>
                        <Td>
                            <Button colorScheme="red" size="md" variant="solid" onClick={() => onDelete.mutate(experiment.id)}>
                                Deletar
                            </Button>
                        </Td>
                    </Tr>
                )
            })}
        </>
    );
}

export default Experiments;