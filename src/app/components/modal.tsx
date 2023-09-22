'use client'
import { Modal as ChakraModal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, ResponsiveValue } from "@chakra-ui/react";

interface Props {
    title?:string,
    isOpen:boolean,
    onClose:() => void,
    size?: ResponsiveValue<(string & {}) | "sm" | "md" | "lg" | "xl" | "2xl" | "xs" | "3xl" | "4xl" | "5xl" | "6xl" | "full">,
    children?:React.ReactNode,
    footer?:React.ReactNode,
}

const Modal = ( {
    title,
    isOpen,
    onClose,
    size,
    children,
    footer,
 } : Props) => {
    return (
        <ChakraModal isOpen={isOpen} onClose={onClose} size={size ? size : 'xl'}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>{title}</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    {children}        
                </ModalBody>
                <ModalFooter>
                    {footer}
                </ModalFooter>
            </ModalContent>
        </ChakraModal>
    );
};

export default Modal;