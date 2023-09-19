'use client'
import { Button, Flex } from "@chakra-ui/react";

const ExperimentSubmission = () => {
    const startExperiment = () => {

    };

    return (
        <Flex justifyContent={'flex-end'}>
            <Button onClick={() => startExperiment()}>Iniciar Experimento</Button>
        </Flex>
    );
};

export default ExperimentSubmission;