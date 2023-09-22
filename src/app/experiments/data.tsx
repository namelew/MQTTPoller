'use client'
import { Tr, Td } from "@chakra-ui/react";
import { IExperiment } from "interfaces/IExperiment";
import { mqttVersions } from "static";

interface Props {
    experiments?:IExperiment[],
}

const Experiments = ( { experiments } : Props) => {
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
                        <Td>Mais</Td>
                        <Td>Excluir</Td>
                    </Tr>
                )
            })}
        </>
    );
}

export default Experiments;