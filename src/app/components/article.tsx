'use client'

import { Divider, Heading, Flex, As } from "@chakra-ui/react";

interface Props {
    children?: React.ReactNode
    title: {
        text: string
        as?: As,
        size?: '2xl' | '3xl' | '4xl' | 'lg' | 'md' | 'sm' | 'xl' | 'xs',
    }
}

const Article = ({ children, title } : Props) => {
    return (
        <>
            <Heading 
                as={title.as ? title.as : 'h2'}
                size={title.size ? title.size : 'xl'}
            >{title.text}</Heading>
            <Divider />
            <Flex textAlign='justify' flexDirection='column' padding={1}>
                { children }
            </Flex>
        </>
    );
};

export default Article;