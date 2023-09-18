'use client'
import { Box, Flex, Text, Button, Link } from "@chakra-ui/react";

export default function Navbar() {
  return (
    <Flex
      as="nav"
      align="center"
      justify="space-between"
      wrap="wrap"
      padding="1.5rem"
      bg="teal.500"
      color="white"
    >
      <Flex align="center" mr={5}>
        <Text fontSize="lg" fontWeight="bold">
            <Link href="/" rel="noopener noreferrer">
                MQTTPoller
            </Link>
        </Text>
      </Flex>

      <Box
        display={{ base: "block", md: "none" }}
        onClick={() => console.log("Toggle Menu")}
      >
        <svg
          fill="white"
          width="12px"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <title>Menu</title>
          <path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z" />
        </svg>
      </Box>

      <Box display={{ base: "none", md: "flex" }} mt={{ base: 4, md: 0 }}>
        <Button bg="transparent" border="1px">
            <Link href="/orquestrators" rel="noopener noreferrer">
                Orquestradores
            </Link>
        </Button>
        <Button bg="transparent" border="1px">
            <Link href="/experiments" rel="noopener noreferrer">
                Experimentos
            </Link>
        </Button>
        <Button bg="transparent" border="1px">
            <Link href="/about" rel="noopener noreferrer">
                Sobre
            </Link>
        </Button>
      </Box>
    </Flex>
  );
}
