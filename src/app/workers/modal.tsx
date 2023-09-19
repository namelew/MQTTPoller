'use client'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, Button, HStack } from "@chakra-ui/react";

interface Props {
    openModal: boolean,
    onClose: () => void,
    selected?: string[],
}

const ExperimentModal = ({ openModal, onClose, selected }: Props) => {
    return (
        <Modal isOpen={openModal} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Parâmetros do Experimento</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    {/* Your modal content goes here */}
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