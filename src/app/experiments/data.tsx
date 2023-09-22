'use client'
import { Tr, Td, Button } from "@chakra-ui/react";
import { useMutation } from "@tanstack/react-query";
import { deleteExperiment } from "consumer";
import { IExperiment } from "interfaces/IExperiment";
import { mqttVersions } from "static";

interface Props {
    experiments?:IExperiment[],
}

const Experiments = ( { experiments } : Props) => {
    const onDelete = useMutation(deleteExperiment, {
        onSuccess: () => alert("Experimento excluido com sucesso"),
        onError: (error) => console.log(error)
    });

    return (
        <>
            {experiments?.map((experiment) => {
                const mqttVersion = mqttVersions.find((version) => version.value === experiment.mqttVersion);
                const mqttVersionName = mqttVersion ? mqttVersion.name : 'Unknown version';

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
                            <Button colorScheme="teal" size="sm" variant="outline" onClick={() => {/* Add your onClick handler here */}}>
                                Detalhar
                            </Button>
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